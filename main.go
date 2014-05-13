package main

import (
	"fmt"
	"os"
)

var commands = []*Command{
	cmdInit,
	cmdList,
	cmdBundle,
	cmdBuild,
	cmdTest,
	cmdHelp,
}

func usage() {
	fmt.Println("dummy")
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name == args[0] && cmd.Runnable() {
			cmd.Run(args)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	usage()
}
