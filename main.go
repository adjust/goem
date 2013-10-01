package main

import (
	"./goem"
	"log"
	"os"
	"strings"
)

var actions = map[string]interface{}{
	"list":   goem.List,
	"bundle": goem.Bundle,
	"build":  goem.Build,
	"test":   goem.Test,
	"help":   goem.Help,
}

func main() {

	// Running goem without action, show help.
	if len(os.Args) <= 1 {
		goem.Help("")
		os.Exit(0)
	}

	// Generate options from all extra command line parameters after 0 and 1.
	var subOption string
	if len(os.Args) > 2 {
		for iter, arg := range os.Args {
			if iter == 0 || iter == 1 {
				continue
			}
			subOption += arg + " "
		}
	}
	subOption = strings.TrimSpace(subOption)

	action := os.Args[1]
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
			v.(func(string))(subOption)
			os.Exit(0)
		} else if k == action && action == "help" {
			v.(func(string))(subOption)
			os.Exit(0)
		}
	}

	goem.Help("")
	log.Println("unknown action")
}
