package goem

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// bundler object is supposed to collect all necessary go third party modules
type Bundler struct {
	config *Config
}

// NewBundler() returns a Bundler object
// on error it exits
// NewBundler sets the GOPATH to ./.go
func NewBundler(config *Config) *Bundler {
	bundler := &Bundler{
		config: config,
	}
	return bundler
}

// bundle() executes all necessary sub methods to fetch or update the current source
// bundle() exits on error
func (self *Bundler) bundle() {
	err := self.makeBase()
	found := false
	if err != nil {
		fmt.Printf("the following error occured while bundling\n")
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
	for _, env := range self.config.Env {
		if getGoEnv() == env.Name {
			self.getPackages(env.Packages)
			found = true
			break
		}
	}
	if !found {
		for _, env := range self.config.Env {
			if "development" == env.Name {
				self.getPackages(env.Packages)
				found = true
				break
			}
		}
	}
}

// build() calls bundle() to ensure updated source()
// build() builds either an a.out binary in the current working dir
// or in the path given with the optional binary name
func (self *Bundler) build(binName string) {
	setGoPath()
	if binName == "" {
		binName = "a.out"
	}

	binName = strings.TrimSpace(binName)
	binName = strings.Replace(binName, "\n", "", -1)

	sourceFiles, err := self.getSourceFiles()
	if err != nil {
		fmt.Printf("while trying to collect source files: " + err.Error())
	}

	goArgs := []string{}
	goArgs = append(goArgs, "build")
	goArgs = append(goArgs, "-o")
	goArgs = append(goArgs, binName)
	goArgs = append(goArgs, sourceFiles...)

	cmd := exec.Command("go", goArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}
}

// makeBase() creates the .go dir and all necessary subdirectories
// makeBase() is called by bundle() to ensure all needed directories exist
// on error makeBase() returns it
func (self *Bundler) makeBase() error {
	goDirs := [3]string{"/src", "/pkg", "/bin"}
	for _, ext := range goDirs {
		err := os.MkdirAll(getGoPath()+ext, 0777)
		if err != nil {
			return fmt.Errorf("while creating needed dirs: " + err.Error())
		}
	}
	return nil
}

// getPackages() checks if there is already the wanted source
// if so it updates the source and sets the head as specified in the Gofile
// otherwise it fetches the source and sets the head
// on error getPackages exits
func (self *Bundler) getPackages(packages []Package) {
	for _, pkg := range packages {
		if self.checkForPath(pkg.Branch) {
			name, err := os.Getwd()
			if err != nil {
				fmt.Printf(err.Error())
			}
			name += "/" + pkg.Branch
			err = os.RemoveAll(getGoPath() + "/src/" + pkg.Name)
			if err != nil {
				fmt.Printf("while trying to remove useless folder: %s", err.Error())
			}
			os.Mkdir(getGoPath()+"/src/", 0777)
			err = os.Symlink(name, getGoPath()+"/src/"+pkg.Name)
			if err != nil {
				fmt.Printf("while trying to set 'self' link: %s\n", err.Error())
			}
			continue
		}
		if pkg.Branch == "self" {
			name, err := os.Getwd()
			if err != nil {
				fmt.Printf(err.Error())
			}
			name += "/" + path.Base(pkg.Name)
			err = os.RemoveAll(getGoPath() + "/src/" + pkg.Name)
			if err != nil {
				fmt.Printf("while trying to remove useless folder: %s", err.Error())
			}
			os.Mkdir(getGoPath()+"/src/"+pkg.Name, 0777)
			err = os.Symlink(name, getGoPath()+"/src/"+pkg.Name+"/"+path.Base(pkg.Name))
			if err != nil {
				fmt.Printf("while trying to set 'self' link: %s\n", err.Error())
			}
			continue
		}
		if !self.sourceExist(pkg) {
			err := self.getSource(pkg)
			if err != nil {
				fmt.Printf("while trying to get the source files: %s\n\n", err.Error())
				os.Exit(1)
			}
		} else {
			err := self.updateSource(pkg)
			if err != nil {
				fmt.Printf("while trying to update the source files: %s\n\n", err.Error())
				os.Exit(1)
			}
		}
		self.setHead(pkg)
	}
}

// checkForPath
func (self *Bundler) checkForPath(path string) bool {
	if path[0] == '/' || path[0] == '.' {
		return true
	}
	return false
}

// sourceExist simpy checks if the source directory already exists
// it returns true if so, false otherwise
func (self *Bundler) sourceExist(pkg Package) bool {
	dir := getGoPath() + "/src/" + pkg.Name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// getSource() downloads the given package via https from github
// on error it tries to remove created dirs and reports on failure to do so
// also it returns the failure
func (self *Bundler) getSource(pkg Package) error {
	err := git.clone(pkg)
	if err != nil {
		delErr := os.RemoveAll(path.Dir(getGoPath() + "/src/" + pkg.Name))
		if delErr != nil {
			fmt.Printf("An error occured while getting the source, but i am unable to tidy up\n")
			fmt.Printf("Please remove %s manually\n\n", path.Dir(getGoPath()+"/src/"+pkg.Name))
		}
		return fmt.Errorf("Could not get source:\n\n" + err.Error())
	}
	return nil
}

// updateSource() tries to call git pull after switching to master branch
// on error it returns the error
func (self *Bundler) updateSource(pkg Package) error {
	err := git.pull(pkg)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	return nil
}

// setHead() tries to set the head according to the Gofile with git checkout
// on error it returns the error
func (self *Bundler) setHead(pkg Package) error {
	err := git.checkout(pkg, "")
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	return nil
}

// getSourceFiles() is called by build() method
// it collects all .go files in the current working dir and returns them as a string
// on error it returns the error
func (self *Bundler) getSourceFiles() ([]string, error) {
	sourceFiles, err := filepath.Glob(self.config.Srcdir + "/*\\.go")
	if err != nil {
		return nil, fmt.Errorf("while trying to get glob filepath: " + err.Error())
	}

	var glob []string
	for _, file := range sourceFiles {
		if ! IsPathDir(file) {
			glob = append(glob, file)
		}
	}
	return glob, nil
}
