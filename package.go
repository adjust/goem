package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Package is a struct to hold a repository name and the desired branch
// Config holds an array of Packages
type Package struct {
	Name   string
	Branch string
	GitUrl string
}

func (self *Package) IsGitHub() bool {
	return strings.HasPrefix(self.Name, "github.com/")
}

func (self *Package) setGitUrl() {
	self.GitUrl = strings.Replace(self.Name, "github.com", "github.com:", -1)
}

func (self *Package) setMirroredGitUrl(mirror string) {
	self.GitUrl = strings.Replace(self.GitUrl, "github.com", mirror, -1)
}

func (self *Package) branchIsPath() bool {
	if self.Branch == "" {
		return false
	}

	if self.Branch[0] == '/' || self.Branch[0] == '.' || self.Branch == "self" {
		return true
	}
	return false
}

func (self *Package) sourceExist() bool {
	dir := getGoPath() + "/src/" + self.Name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func (self *Package) getSource() {
	var err error
	if self.IsGitHub() {
		err = git.clone(self)
	} else {
		var out []byte
		if out, err = exec.Command("go", "get", "-d", self.Name).CombinedOutput(); err != nil {
			err = errors.New(string(out))
		}
	}

	if err != nil {
		delErr := os.RemoveAll(path.Dir(getGoPath() + "/src/" + self.Name))
		if delErr != nil {
			fmt.Printf("An error occured while getting the source, but I am unable to tidy up\n")
			fmt.Printf("Please remove %s manually\n\n", path.Dir(getGoPath()+"/src/"+self.Name))
		}
		stderrAndExit(err)
	}
}

func (self *Package) updateSource() {
	var err error
	if self.IsGitHub() {
		err = git.pull(self)
	} else {
		var out []byte
		if out, err = exec.Command("go", "get", "-d", "-u", self.Name).CombinedOutput(); err != nil {
			err = errors.New(string(out))
		}
	}

	if err != nil {
		stderrAndExit(err)
	}
}

func (self *Package) setHead() {
	err := git.checkout(self, "")
	if err != nil {
		stderrAndExit(err)
	}
}

func (self *Package) createSymlink() {
	name, err := os.Getwd()
	if err != nil {
		stderrAndExit(err)
	}
	name += "/" + self.Branch
	err = os.RemoveAll(getGoPath() + "/src/" + self.Name)
	if err != nil {
		stderrAndExit(err)
	}
	os.MkdirAll(path.Dir(getGoPath()+"/src/"+self.Name), 0777)
	err = os.Symlink(name, getGoPath()+"/src/"+self.Name)
	if err != nil {
		stderrAndExit(err)
	}
}

type Packages []*Package

func (self Packages) Len() int {
	return len(self)
}

func (self Packages) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type ByName struct{ Packages }

func (self *ByName) Less(i, j int) bool {
	return self.Packages[i].Name < self.Packages[j].Name
}
