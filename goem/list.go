package goem

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Lister struct {
	goPath string
}

func NewList() *Lister {
	goPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	} else {
		goPath += "/.go/src/"
	}
	list := &Lister{
		goPath: goPath,
	}
	return list
}

func (self *Lister) list() {
	self.dirRead(0, self.goPath)
}

func (self *Lister) dirRead(called int, path string) {
	called++
	if called == 4 {
		return
	}
	dirGlob, err := filepath.Glob(path + "/*")
	if err != nil {
		log.Println(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	for _, dir := range dirGlob {
		self.dirRead(called, dir)
		if called == 3 {
			log.Println(strings.Replace(dir, cwd+"/.go/src/", "", -1))
		}
	}
}
