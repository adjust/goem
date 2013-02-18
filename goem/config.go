package goem

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// GoPkg is a struct that relates Configs with there Package name
type GoPkg struct {
	name   string
	config *Config
}

// Config is a struct that holds the information from the Gofile
// config has an array of environments which have a name and an array of
// third party packages
type Config struct {
	Env []Env
}

type Env struct {
	Name     string
	Packages []Package
}

// NewConfig() calls parse()
// NewConfig() simply returns an initialized Config object
func NewConfig() *Config {
	config := &Config{}
	config.parse("")
	return config
}

// NewLockConfig()
// same as NewConfig(), but initializes with Gofile.lock
func NewLockConfig() *Config {
	config := &Config{}
	config.parse("Gofile.lock")
	return config
}

// parse() reads the Gofile, which is expected to be in the current
// working dir. After reading the Gofile the content is unmarshaled into
// the Config object
// parse exits on error
func (self *Config) parse(gofile string) {
	if gofile == "" {
		gofile = "./Gofile"
	}
	configData, err := ioutil.ReadFile(gofile)
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
