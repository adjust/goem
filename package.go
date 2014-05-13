package main

import (
	"os"
)

// Package is a struct to hold a repository name and the desired branch
// Config holds an array of Packages
type Package struct {
	Name   string
	Branch string
}

func (self *Package) BranchIsPath() bool {
	if self.Branch[0] == '/' || self.Branch[0] == '.' {
		return true
	}
	return false
}

func (self *Package) SourceExist() bool {
	dir := getGoPath() + "/src/" + self.Name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

type Packages []Package

func (self *Packages) Len() int {
	return len(*self)
}

func (self Packages) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type ByName struct{ Packages }

func (self *ByName) Less(i, j int) bool {
	return self.Packages[i].Name < self.Packages[j].Name
}
