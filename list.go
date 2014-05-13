package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var srcPath string = getGoPath() + "/src/"

var cmdList = &Command{
	Run:  list,
	Name: "list",
}

// list() prints all installed go extensions to stdout
// it does so by calling dirRead()
// on error it prints it and exits
func list(args []string) {
	results, err := dirRead(0, srcPath, nil)
	if err != nil {
		stderrAndExit(err)
	}
	for _, result := range results {
		fmt.Printf("%s\n", result)
	}
}

// dirRead crawls the current GOPATH and returns an array which holds all packages
// if an error occurs it returns it, nil otherwise
func dirRead(called int, path string, result []string) ([]string, error) {
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
		result, _ = dirRead(called, dir, result)
		if called == 3 {
			temp := make([]string, len(result)+1)
			for i, v := range result {
				temp[i] = v
			}
			result = temp
			goSrcDir := filepath.FromSlash(cwd + "/.go/src/")
			result[len(result)-1] = strings.Replace(dir, goSrcDir, "", -1)
		}
	}
	return result, nil
}
