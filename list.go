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

func list(args []string) {
	results := dirRead(0, srcPath, nil)
	for _, result := range results {
		fmt.Printf("%s\n", result)
	}
}

func dirRead(called int, path string, result []string) []string {
	if result == nil {
		result = []string{}
	}
	called++
	if called == 4 {
		return result
	}
	dirGlob, err := filepath.Glob(path + "/*")
	if err != nil {
		stderrAndExit(err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		stderrAndExit(err)
	}
	for _, dir := range dirGlob {
		result = dirRead(called, dir, result)
		if called == 3 {
			goSrcDir := filepath.FromSlash(cwd + "/.go/src/")
			result = append(result, strings.Replace(dir, goSrcDir, "", -1))
		}
	}
	return result
}
