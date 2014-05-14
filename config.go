package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var cmdInit = &Command{
	Run:  initConfig,
	Name: "init",
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

func initConfig(args []string) {
	gofile := "./Gofile"
	if fileExists(gofile) {
		stderrAndExit(fmt.Errorf("Gofile already exists"))
	}
	devEnv := Env{"development", []Package{}}
	config := &Config{[]Env{devEnv}, "./test", ""}
	config.write("./Gofile")
}

func (self *Config) parse(gofile string) {
	if gofile == "" {
		gofile = "./Gofile"
	}
	configData, err := ioutil.ReadFile(gofile)
	if err != nil {
		stderrAndExit(err)
	}
	err = json.Unmarshal(configData, self)
	if err != nil {
		stderrAndExit(err)
	}
}

func (self *Config) write(gofile string) {
	j, err := json.MarshalIndent(self, "", "\t")
	if err != nil {
		stderrAndExit(err)
	}
	err = ioutil.WriteFile(gofile, j, 0777)
	if err != nil {
		stderrAndExit(err)
	}
}
