package goem

import (
	"fmt"
	"os"
)

func IsPathDir(path string) bool {
	osFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("while isDir() open: " + err.Error())
		return false
	}
	defer osFile.Close()

	statFile, err := osFile.Stat()
	if err != nil {
		fmt.Printf("while isDir() stat: " + err.Error())
		return false
	}

	if statFile.Mode().IsDir() {
		return true
	}
	return false
}
