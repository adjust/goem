package goem

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Bundler struct {
	goPath string
	goEnv  string
	config *Config
}

func NewBundler(config *Config) *Bundler {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		goEnv = "development"
	}
	goPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
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

func (self *Bundler) bundle() {
	self.makeBase()
	for _, env := range self.config.Env {
		if self.goEnv == env.Name {
			self.getPackages(env.Packages)
		}
	}
}

func (self *Bundler) build(binName string) {
	self.bundle()
	cwd, _ := os.Getwd()
	os.Setenv("GOPATH", cwd+"/.go")
	log.Println(os.Getenv("GOPATH"))
	if binName == "" {
		binName = "a.out"
	}
	out, err := exec.Command(
		"/usr/bin/go",
		"build",
		"-o",
		binName,
		self.getSourceFiles(),
	)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func (self *Bundler) makeBase() {
	goDirs := [3]string{"/src", "/pkg", "/bin"}
	for _, ext := range goDirs {
		err := os.MkdirAll(self.goPath+ext, 0777)
		if err != nil {
			log.Println(err)
		}
	}
}

func (self *Bundler) getPackages(packages []Package) {
	for _, pkg := range packages {
		if !self.sourceExist(pkg) {
			self.getSource(pkg)
		} else {
			self.updateSource(pkg)
		}
		self.setHead(pkg)
	}
}

func (self *Bundler) sourceExist(pkg Package) bool {
	dir := self.goPath + "/src/" + pkg.Name
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	}
	return true
}

func (self *Bundler) getSource(pkg Package) {
	cmd := exec.Command(
		"git",
		"clone",
		"https://"+pkg.Name+".git",
		self.goPath+"/src/"+pkg.Name,
	)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func (self *Bundler) updateSource(pkg Package) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	dir := self.goPath + "/src/" + pkg.Name
	os.Chdir(dir)
	cmd := exec.Command(
		"git",
		"checkout",
		"master",
	)
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
	cmd = exec.Command(
		"git",
		"pull",
	)
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
	os.Chdir(currentDir)
}

func (self *Bundler) setHead(pkg Package) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	dir := self.goPath + "/src/" + pkg.Name
	os.Chdir(dir)
	cmd := exec.Command(
		"git",
		"checkout",
		pkg.Branch,
	)
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
	os.Chdir(currentDir)
}

func (self *Bundler) getSourceFiles() string {
	cwd, err := os.Getwd()
	var glob string
	if err != nil {
		log.Println(err)
	}
	sourceFiles, err := filepath.Glob(cwd + "/*\\.go")
	if err != nil {
		log.Println(err)
	}
	regex := regexp.MustCompile("^\\.go")
	for _, file := range sourceFiles {
		base := path.Base(file)
		if !regex.MatchString(base) {
			glob += base + " "
		}
	}
	log.Println(strings.TrimSpace(glob))
	return strings.TrimSpace(glob)
}
