package goem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// List is a struct that simply holds the current GOPATH
type Lister struct {
	srcPath string
}

// NewList() returns a new ListObject
// GOPATH is set to "./go/src/"
// on error it exits
func NewList() *Lister {
	list := &Lister{
		srcPath: getGoPath() + "/src/",
	}
	return list
}

// list() prints all installed go extensions to stdout
// it does so by calling dirRead()
// on error it prints it and exits
func (self *Lister) list() {
	results, err := self.dirRead(0, self.srcPath, nil)
	if err != nil {
		fmt.Printf("while dirRead(): " + err.Error())
		os.Exit(1)
	}
	for _, result := range results {
		fmt.Printf("%s\n", result)
	}
}

// dirRead crawls the current GOPATH and returns an array which holds all packages
// if an error occurs it returns it, nil otherwise
func (self *Lister) dirRead(called int, path string, result []string) ([]string, error) {
	if result == nil {
		result = make([]string, 1)
	}
	called++
	if called == 4 {
		return result, nil
	}
	dirGlob, err := filepath.Glob(path + "/*")
	if err != nil {
		return nil, fmt.Errorf("while calling Glob on %s: %s\n", path, err.Error())
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("while trying to get current working dir: " + err.Error())
	}
	for _, dir := range dirGlob {
		result, _ = self.dirRead(called, dir, result)
		if called == 3 {
			temp := make([]string, len(result)+1)
			for i, v := range result {
				temp[i] = v
			}
			result = temp
			goSrcDir := filepath.FromSlash(cwd+"/.go/src/")
			result[len(result)-1] = strings.Replace(dir, goSrcDir, "", -1)
		}
	}
	return result, nil
}
