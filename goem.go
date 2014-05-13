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
		fmt.Printf("Could not construct Bundler Object: %s\n", err.Error())
		os.Exit(1)
	} else {
		goPath += "/.go/"
	}
	return goPath
}

func setGoPath() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("while trying to get working dir: " + err.Error())
	}

	os.Setenv("GOPATH", cwd+"/.go")
}

func stderrAndExit(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
