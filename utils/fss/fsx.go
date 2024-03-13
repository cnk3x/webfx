package fss

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func SaveFile(src io.Reader, dstPath string, overwrite bool) (hash string, err error) {
	if err = os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return
	}

	fFlags := os.O_RDWR | os.O_CREATE
	if overwrite {
		fFlags |= os.O_TRUNC
	} else {
		fFlags |= os.O_EXCL
	}

	var dst *os.File
	if dst, err = os.OpenFile(dstPath, fFlags, 0o666); err != nil {
		return
	}
	defer dst.Close()

	m := md5.New()
	if _, err = io.Copy(io.MultiWriter(dst, m), src); err != nil {
		return
	}
	hash = hex.EncodeToString(m.Sum(nil))
	return
}

func MakeDirs(dir string) string {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return dir
}
