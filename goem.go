package main

import (
	"fmt"
	"os"
)

func getGoEnv() string {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		goEnv = "development"
	}
	return goEnv
}

func getGoPath() string {
	goPath, err := os.Getwd()
	if err != nil {
		stderrAndExit(err)
	}
	goPath += "/.go/"
	return goPath
}

func setGoPath() {
	cwd, err := os.Getwd()
	if err != nil {
		stderrAndExit(err)
	}

	os.Setenv("GOPATH", cwd+"/.go")
}

func stderrAndExit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
