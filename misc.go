package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/fs"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Slash(text string) (string, string) {
	res := strings.SplitN(text, "/", 2)
	if len(res) == 1 {
		return res[0], ""
	}
	return res[0], res[1]
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

func Dump(path string, what interface{}) error {
	if b, err := json.MarshalIndent(what, "", "\t"); err != nil {
		return Errorc(err)
	} else if err := os.WriteFile(path, b, 0644); err != nil {
		return Errorc(err)
	}

	return nil
}

func ReadDirByDate(dirname string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, Errorc(err)
	}

	// Create a slice of struct that includes os.DirEntry and modification time
	type entryWithTime struct {
		entry   fs.DirEntry
		modTime time.Time
	}

	entryWrappers := []entryWithTime{}

	for _, entry := range entries {
		info, err := entry.Info()
		if errors.Is(err, fs.ErrNotExist) {
			continue // file got deleted while we were working
		} else if err != nil {
			return nil, Errorc(err)
		}
		entryWrappers = append(entryWrappers, entryWithTime{entry: entry, modTime: info.ModTime()})
	}

	// Sort the entries by modification time
	sort.Slice(entryWrappers, func(i, j int) bool {
		return entryWrappers[i].modTime.Before(entryWrappers[j].modTime)
	})

	// Extract sorted os.DirEntry from the wrappers
	sortedEntries := make([]fs.DirEntry, len(entries))
	for i, wrapper := range entryWrappers {
		sortedEntries[i] = wrapper.entry
	}

	return sortedEntries, nil
}
