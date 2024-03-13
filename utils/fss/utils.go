package fss

import "os"

func IsFile(file string) bool {
	f, _ := os.Stat(file)
	return f != nil && f.Mode().IsRegular()
}

func IsDIR(file string) bool {
	f, _ := os.Stat(file)
	return f != nil && f.IsDir()
}
