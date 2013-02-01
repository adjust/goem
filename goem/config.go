package goem

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Env []struct {
		Name     string
		Packages []Package
	}
}

type Package struct {
	Name   string
	Branch string
}

func NewConfig() *Config {
	config := &Config{}
	config.parse()
	return config
}

func (self *Config) parse() {
	configData, err := ioutil.ReadFile("./Gofile")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(configData, self)
	if err != nil {
		log.Println(err)
	}
}
