package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		helpHelp()
		os.Exit(1)
	}
	for _, cmd := range commands {
		if cmd.Name == args[0] && cmd.Runnable() {
			cmd.Run(args[1:])
			return
		}
	}
	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	helpHelp()
}
