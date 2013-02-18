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
}

func main() {
	var subOption string
	action := os.Args[1]
	if len(os.Args) > 2 {
		subOption = os.Args[2]
	}
	for k, v := range actions {
		if k == action && action == "list" {
			v.(func())()
			os.Exit(0)
		} else if k == action && action == "bundle" {
			v.(func())()
			os.Exit(0)
		} else if k == action && action == "build" {
			v.(func(string))(subOption)
			os.Exit(0)
		}
	}
	log.Println("unknown action")
}
