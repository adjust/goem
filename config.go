package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var cmdInit = &Command{
	Run:  InitConfig,
	Name: "init",
}

type GoPkg struct {
	name   string
	config *Config
}

type Config struct {
	Env     []Env  `json:"env"`
	Testdir string `json:"testdir"`
	Mirror  string `json:"mirror"`
}

type Env struct {
	Name     string    `json:"name"`
	Packages []Package `json:"packages"`
}

var config *Config = &Config{}

func InitConfig(args []string) {
	gofile := "./Gofile"

	dev_env := Env{"development", []Package{}}
	config := &Config{[]Env{dev_env}, "./test", ""}

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
}

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
