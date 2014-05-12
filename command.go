package main

import (
	"flag"
	"fmt"
	"strings"
)

type Command struct {
	// args does not include the command name
	Run  func(args []string)
	Flag flag.FlagSet

	Usage string // first word is the command name
	Short string // `redismq help` output
	Long  string // `redismq help cmd` output
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
