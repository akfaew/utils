package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
