package goem

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// bundler object is supposed to collect all necessary go third party modules
type Bundler struct {
	goPath string
	goEnv  string
	config *Config
}

// NewBundler() returns a Bundler object
// on error it exits
// NewBundler sets the GOPATH to ./.go
func NewBundler(config *Config) *Bundler {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		goEnv = "development"
	}
	goPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("Could not construct Bundler Object: %s\n", err.Error())
		os.Exit(1)
	} else {
		goPath += "/.go"
	}
	bundler := &Bundler{
		goPath: goPath,
		config: config,
		goEnv:  goEnv,
	}
	return bundler
}

// bundle() executes all necessary sub methods to fetch or update the current source
// bundle() exits on error
func (self *Bundler) bundle() {
	err := self.makeBase()
	if err != nil {
		fmt.Printf("the following error occured while bundling\n")
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
	for _, env := range self.config.Env {
		if self.goEnv == env.Name {
			self.getPackages(env.Packages)
		}
	}
}

// build() calls bundle() to ensure updated source()
// build() builds either an a.out binary in the current working dir
// or in the path given with the optional binary name
func (self *Bundler) build(binName string) {
	self.bundle()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("while trying to get working dir: " + err.Error())
	}

	os.Setenv("GOPATH", cwd+"/.go")
	if binName == "" {
		binName = "a.out"
	}

	sourceFiles, err := self.getSourceFiles()
	if err != nil {
		fmt.Printf("while trying to collect source files: " + err.Error())
	}

	cmd := exec.Command(
		"/usr/bin/go",
		"build",
		"-o",
		binName,
		sourceFiles,
	)
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
		err := os.MkdirAll(self.goPath+ext, 0777)
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

// sourceExist simpy checks if the source directory already exists
// it returns true if so, false otherwise
func (self *Bundler) sourceExist(pkg Package) bool {
	dir := self.goPath + "/src/" + pkg.Name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

// getSource() downloads the given package via https from github
// on error it tries to remove created dirs and reports on failure to do so
// also it returns the failure
func (self *Bundler) getSource(pkg Package) error {
	cmd := exec.Command(
		"git",
		"clone",
		"https://"+pkg.Name+".git",
		self.goPath+"/src/"+pkg.Name,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		delErr := os.RemoveAll(path.Dir(self.goPath + "/src/" + pkg.Name))
		if delErr != nil {
			fmt.Printf("An error occured while getting the source, but i am unable to tidy up\n")
			fmt.Printf("Please remove %s manually\n\n", path.Dir(self.goPath+"/src/"+pkg.Name))
		}
		return fmt.Errorf("Could not get source:\n\n%s\n"+err.Error(), out)
	}
	return nil
}

// updateSource() tries to call git pull after switching to master branch
// on error it returns the error
func (self *Bundler) updateSource(pkg Package) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("while trying to get working dir: " + err.Error())
	}
	dir := self.goPath + "/src/" + pkg.Name
	os.Chdir(dir)
	cmd := exec.Command(
		"git",
		"checkout",
		"master",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Could not switch branch:\n\n%s\n"+err.Error(), out)
	}
	cmd = exec.Command(
		"git",
		"pull",
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Could not update source:\n\n%s\n"+err.Error(), out)
	}
	os.Chdir(currentDir)
	return nil
}

// setHead() tries to set the head according to the Gofile with git checkout
// on error it returns the error
func (self *Bundler) setHead(pkg Package) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("while trying to get working dir: " + err.Error())
	}
	dir := self.goPath + "/src/" + pkg.Name
	os.Chdir(dir)
	cmd := exec.Command(
		"git",
		"checkout",
		pkg.Branch,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Could not switch head:\n\n%s\n"+err.Error(), out)
	}
	os.Chdir(currentDir)
	return nil
}

// getSourceFiles() is called by build() method
// it collects all .go files in the current working dir and returns them as a string
// on error it returns the error
func (self *Bundler) getSourceFiles() (string, error) {
	cwd, err := os.Getwd()
	var glob string
	if err != nil {
		return "", fmt.Errorf("while trying to get working dir: " + err.Error())
	}
	sourceFiles, err := filepath.Glob(cwd + "/*\\.go")
	if err != nil {
		return "", fmt.Errorf("while trying to get glob filepath: " + err.Error())
	}
	regex := regexp.MustCompile("^\\.go")
	for _, file := range sourceFiles {
		base := path.Base(file)
		if !regex.MatchString(base) {
			glob += base + " "
		}
	}
	return strings.TrimSpace(glob), nil
}
