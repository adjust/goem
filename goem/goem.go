package goem

import (
	"os"
)

func List() {
	lister := NewList()
	lister.list()
	os.Exit(0)
}

func Bundle() {
	config := NewConfig()
	bundler := NewBundler(config)
	bundler.bundle()
	os.Exit(0)
}

func Build(binName string) {
	config := NewConfig()
	bundler := NewBundler(config)
	bundler.build(binName)
	os.Exit(0)
}
