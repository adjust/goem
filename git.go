package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Git struct{}

var git Git

func (self *Git) clone(pkg *Package) error {
	cmd := exec.Command(
		"git",
		"clone",
		"git@"+gitUrl+".git",
		getGoPath()+"/src/"+pkg.Name,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone:\n%s\n%s\n", err.Error(), out)
	}
	return nil
}

func (self *Git) checkout(pkg *Package, branch string) error {
	if branch == "" {
		branch = self.checkBranch(pkg.Branch)
	}

	oldDir := self.dirSwap(pkg, "")

	cmd := exec.Command(
		"git",
		"checkout",
		branch,
	)
	out, err := cmd.CombinedOutput()

	self.dirSwap(pkg, oldDir)

	if err != nil {
		return fmt.Errorf("Could not switch branch:\n\n%s\n"+err.Error(), out)
	}
	return nil
}

func (self *Git) pull(pkg *Package) error {
	err := self.checkout(pkg, "master")
	if err != nil {
		return fmt.Errorf("git pull: %s", err.Error())
	}

	oldDir := self.dirSwap(pkg, "")

	cmd := exec.Command(
		"git",
		"pull",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"Could not update source:%s\n%s\n%s\n",
			pkg.Name,
			err.Error(),
			out,
		)
	}

	self.dirSwap(pkg, oldDir)

	return nil
}

func (self *Git) log(pkg *Package, format string) ([]string, error) {
	if format == "" {
		format = "%H"
	}
	self.checkout(pkg, "master")
	oldDir := self.dirSwap(pkg, "")
	cmd := exec.Command(
		"/usr/bin/git",
		"log",
		"--decorate",
		"--pretty=format:"+format,
	)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("while trying to set out pipe: %s\n", err.Error())
	}

	if err = cmd.Start(); err != nil {
		return nil, fmt.Errorf("while trying to get git log: %s \n", err.Error())
	}

	data, err := ioutil.ReadAll(out)
	if err != nil {
		return nil, fmt.Errorf("while trying to read from out pipe: %s\n", err.Error())
	}

	self.dirSwap(pkg, oldDir)
	self.checkout(pkg, "")
	return strings.Split(string(data), "\n"), nil
}

func (self *Git) refNameToCommit(pkg *Package) (string, error) {
	gitlog, err := self.log(pkg, "%H %d")
	if err != nil {
		return "", fmt.Errorf("refNameToCommit: %s", err.Error())
	}
	regex := regexp.MustCompile(self.checkBranch(pkg.Branch))
	for _, line := range gitlog {
		if regex.MatchString(line) {
			/*
				t1 := strings.Split(line, " ")
				t2 := strings.Replace(t1[1], "(", "", -1)
				t2 = strings.Replace(t2, ")", "", -1)
				t3 = strings.Split(t2, ",")
				t4 := len(t3)
				name := t3[t4-1]
				name = string.TrimSpace(name)
			*/
			if pkg.Branch[0] == '<' || pkg.Branch[0] == '>' {
				return pkg.Branch[0:1] + strings.Split(line, " ")[0], nil
			}
			return strings.Split(line, " ")[0], nil
		}
	}
	return "", nil
}

func (self *Git) dirSwap(pkg *Package, dir string) string {
	oldDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not switch dir: %s\n", err.Error())
	}
	if dir == "" {
		dir = getGoPath() + "/src/" + pkg.Name
	}
	os.Chdir(dir)
	return oldDir
}

func (self *Git) checkBranch(branch string) string {
	if branch[0] == '>' || branch[0] == '<' {
		return branch[1:]
	}
	return branch
}
