package goem

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config is a struct that holds the information from the Gofile
// config has an array of environments which have a name and an array of
// third party packages
type Config struct {
	Env []struct {
		Name     string
		Packages []Package
	}
}

// Package is a struct to hold a repository name and the desired branch
// Config holds an array of Packages
type Package struct {
	Name   string
	Branch string
}

// NewConfig() calls parse()
// NewConfig() simply returns an initialized Config object
func NewConfig() *Config {
	config := &Config{}
	config.parse()
	return config
}

// parse() reads the Gofile, which is expected to be in the current
// working dir. After reading the Gofile the content is unmarshaled into
// the Config object
// parse exits on error
func (self *Config) parse() {
	configData, err := ioutil.ReadFile("./Gofile")
	if err != nil {
		fmt.Printf("while trying to read Gofile: " + err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(configData, self)
	if err != nil {
		fmt.Printf("while trying to unmarshal Gofile: " + err.Error())
		os.Exit(1)
	}
}
