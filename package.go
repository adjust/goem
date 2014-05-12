package main

// Package is a struct to hold a repository name and the desired branch
// Config holds an array of Packages
type Package struct {
	Name   string
	Branch string
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
