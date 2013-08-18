package goem

import (
	"fmt"
	"os"
)

var git Git

// goem just holds a few fuctions, which create the needed objects
// all these functions exit after being called
// all goem functions only work in the root directory of the project

// List() is supposed to list all go extensions installed is the current project
func List() {
	lister := NewList()
	lister.list()
}

// Bundle() is supposed to get all packages specified in the Gofile
// if the package is already installed Bundle() will update the git repo of the package
// afterwards it will set the head of the package according to the Gofile entry
func Bundle(subOption string) {
	config := NewConfig()
	bundler := NewBundler(config)
	bundler.bundle()
	dep := NewDepCheck(config, subOption)
	dep.Start()
	config = NewLockConfig()
	bundler = NewBundler(config)
	bundler.bundle()
}

// Build() builds the binary file
// if no output file is set it will create a.out in the current working dir
func Build(binName string) {
	config := NewLockConfig()
	bundler := NewBundler(config)
	bundler.build(binName)
}

func Test(testDir string) {
	config := NewConfig()
	if !test(config, testDir) {
		os.Exit(1)
	}
}

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
