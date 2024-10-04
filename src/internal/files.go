package internal

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

type HashAlgo int

const (
	MD5    HashAlgo = iota // MD5 = 0
	SHA256                 // SHA256 = 1
)

func ScanPath(db *SqliteDB, session int64, path string, hashAlgo HashAlgo) error {
	files, err := filepath.Glob(path + "/**/*")
	if err != nil {
		return err
	}
	for _, file := range files {
		hash, err := hashFile(file, MD5)
		if err != nil {
			return err
		}
		id, isDuplicate, err := db.AddHash(session, hashAlgo, hash)
		if err != nil {
			return err
		}
		_, err = db.AddFile(session, path, file, id)
		if err != nil {
			return err
		}
		if isDuplicate {
			db.logger.Trace("Duplicate hash", hash)
		}

	}

	return nil
}

func hashFile(path string, algo HashAlgo) (string, error) {
	//Calculate md5
	var hash hash.Hash
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	switch algo {
	case MD5:
		hash = md5.New()
	case SHA256:
		hash = sha256.New()
	}
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil
}
