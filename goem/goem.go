package goem

import (
	"os"
)

// goem just holds a few fuctions, which create the needed objects
// all these functions exit after being called
// all goem functions only work in the root directory of the project

// List() is supposed to list all go extensions installed is the current project
func List() {
	lister := NewList()
	lister.list()
	os.Exit(0)
}

// Bundle() is supposed to get all packages specified in the Gofile
// if the package is already installed Bundle() will update the git repo of the package
// afterwards it will set the head of the package according to the Gofile entry
func Bundle() {
	config := NewConfig()
	bundler := NewBundler(config)
	bundler.bundle()
	os.Exit(0)
}

// Build() builds the binary file
// if no output file is set it will create a.out in the current working dir
// Build() calls Bundle() before building the binary so calling Bundle() is useless
func Build(binName string) {
	config := NewConfig()
	bundler := NewBundler(config)
	bundler.build(binName)
	os.Exit(0)
}
