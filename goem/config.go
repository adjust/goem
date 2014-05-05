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
	Env     []Env  `json:"env"`
	Testdir string `json:"testdir"`
}

type Env struct {
	Name     string    `json:"name"`
	Packages []Package `json:"packages"`
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

func InitConfig() (*Config) {
	gofile := "./Gofile"

	dev_env := Env{"development", []Package{}}
	config := &Config{[]Env{dev_env}, "./test"}

	f, err := os.OpenFile(gofile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if os.IsExist(err) {
		fmt.Printf("There is already a Gofile in current directory. Skipping.\n")
		os.Exit(0)
	} else if err != nil {
		fmt.Printf("Failed to create Gofile: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	configData, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		fmt.Printf("Failed to write to Gofile: %s\n", err.Error())
	}

	_, err = f.Write(append(configData, '\n'))
	if err != nil {
		fmt.Printf("Failed to write to Gofile: %s\n", err.Error())
	}

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
