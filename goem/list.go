package goem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// List is a struct that simply holds the current GOPATH
type Lister struct {
	goPath string
}

// NewList() returns a new ListObject
// GOPATH is set to "./go/src/"
// on error it exits
func NewList() *Lister {
	goPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("while trying to get current working dir: " + err.Error())
		os.Exit(1)
	} else {
		goPath += "/.go/src/"
	}
	list := &Lister{
		goPath: goPath,
	}
	return list
}

// list() prints all installed go extensions to stdout
// it does so by calling dirRead()
// on error it prints it and exits
func (self *Lister) list() {
	err := self.dirRead(0, self.goPath)
	if err != nil {
		fmt.Printf("while dirRead(): " + err.Error())
		os.Exit(1)
	}
}

// dirRead crawls the current GOPATH and prints all directories to stdout
// if an error occurs it returns it, nil otherwise
func (self *Lister) dirRead(called int, path string) error {
	called++
	if called == 4 {
		return
	}
	dirGlob, err := filepath.Glob(path + "/*")
	if err != nil {
		return fmt.Errorf("while calling Glob on %s: %s\n", path, err.Error())
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("while trying to get current working dir: " + err.Error())
	}
	for _, dir := range dirGlob {
		self.dirRead(called, dir)
		if called == 3 {
			fmt.Println(strings.Replace(dir, cwd+"/.go/src/", "", -1))
		}
	}
	return nil
}
