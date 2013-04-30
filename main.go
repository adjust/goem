package main

import (
	"./goem"
	"log"
	"os"
)

var actions = map[string]interface{}{
	"list":   goem.List,
	"bundle": goem.Bundle,
	"build":  goem.Build,
	"test":   goem.Test,
	"help":   goem.Help,
}

func main() {
	var subOption string
	action := os.Args[1]
	if len(os.Args) > 2 {
		for iter, arg := range os.Args {
			if iter == 0 || iter == 1 {
				continue
			}
			subOption += arg + " "
		}
	}
	for k, v := range actions {
		if k == action && action == "list" {
			v.(func())()
			os.Exit(0)
		} else if k == action && action == "bundle" {
			v.(func(string))(subOption)
			os.Exit(0)
		} else if k == action && action == "build" {
			v.(func(string))(subOption)
			os.Exit(0)
		} else if k == action && action == "test" {
			v.(func())()
			os.Exit(0)
		} else if k == action && action == "help" {
			v.(func(string))(subOption)
			os.Exit(0)
		}
	}
	goem.Help("")
	log.Println("unknown action")
}
