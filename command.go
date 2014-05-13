package main

import (
	"fmt"
	"strings"
)

type Command struct {
	// args does not include the command name
	Run   func(args []string)
	Usage string
	Short string
	Long  string
	Name  string
}

func (c *Command) printUsage() {
	if c.Runnable() {
		fmt.Printf("Usage: redismq-cli %s\n\n", c.Usage)
	}
	fmt.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}
