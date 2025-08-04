package utils

import (
	"cmp"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/fs"
	rand "math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func Slash(text string) (string, string) {
	before, after, found := strings.Cut(text, "/")
	if !found {
		return before, ""
	}
	return before, after
}

// A simple sum for naming fixture files in tests, e.g. based on an URL.
func Crc32(str string) string {
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(str), crc32.IEEETable))
}

// Generate a random string to append to emails, so that Gmail doesn't clump
// them into "conversations".
func RandEmail() string {
	randomValue := strconv.Itoa(rand.Int())

	hasher := md5.New()
	hasher.Write([]byte(randomValue))
	md5Hash := hex.EncodeToString(hasher.Sum(nil))

	return md5Hash[:len(md5Hash)/2]
}

func Dump(path string, what any) error {
	if b, err := json.MarshalIndent(what, "", "\t"); err != nil {
		return Errorc(err)
	} else if err := os.WriteFile(path, b, 0644); err != nil {
		return Errorc(err)
	}

	return nil
}

// ReadDirByDate reads the directory and returns entries sorted by modification time (ascending: oldest first).
func ReadDirByDate(dirname string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, Errorc(err)
	}

	type entryWithTime struct {
		entry   fs.DirEntry
		modTime time.Time
	}

	entryWrappers := make([]entryWithTime, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if errors.Is(err, fs.ErrNotExist) {
			continue // file got deleted while we were working
		} else if err != nil {
			return nil, Errorc(err)
		}
		entryWrappers = append(entryWrappers, entryWithTime{entry: entry, modTime: info.ModTime()})
	}

	// Sort by modification time (ascending)
	slices.SortFunc(entryWrappers, func(a, b entryWithTime) int {
		return a.modTime.Compare(b.modTime)
	})

	// Extract sorted entries
	sortedEntries := make([]fs.DirEntry, len(entryWrappers))
	for i, wrapper := range entryWrappers {
		sortedEntries[i] = wrapper.entry
	}

	return sortedEntries, nil
}

func Uniq[T cmp.Ordered](s []T) []T {
	return slices.Compact(slices.Sorted(slices.Values(s)))
}
