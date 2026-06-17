package fileops

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)

func VerifyHash(filePath string, hashRef string) (bool, error) {
	var hashType string
	var checksum string

	parts := strings.SplitN(hashRef, ":", 2)
	if len(parts) == 1 {
		hashType = "sha256"
		checksum = strings.ToLower(parts[0])
	} else {
		hashType = strings.ToLower(parts[0])
		checksum = strings.ToLower(parts[1])
	}

	fileHandle, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer func(fileHandle *os.File) {
		err := fileHandle.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(fileHandle)

	var hasher hash.Hash
	switch hashType {
	case "sha256":
		hasher = sha256.New()
	case "sha512":
		hasher = sha512.New()
	case "sha1":
		hasher = sha1.New()
	case "md5":
		hasher = md5.New()
	default:
		return false, fmt.Errorf("unknown hash type: %s", hashType)
	}

	if _, err := io.Copy(hasher, fileHandle); err != nil {
		return false, err
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))
	return fileHash == checksum, nil
}
