package main

import (
	"fmt"
	"os"
	"os/exec"
)

var cmdTest = &Command{
	Run:  test,
	Name: "test",
}

func test(args []string) {
	config.parse("")
	var testDir string
	if len(args) > 0 {
		testDir = args[len(args)-1]
	}
	setGoPath()
	buildPackages := false
	if testDir == "" {
		testDir = config.Testdir
	}
	for _, arg := range args {
		if arg == "-i" {
			buildPackages = true
		}
	}
	os.Chdir(testDir)

	testCommand := ""
	if buildPackages {
		testCommand = "-i"
	}

	execTest := exec.Command(
		"go",
		"test",
		testCommand,
	)
	out, err := execTest.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		os.Exit(1)
	} else {
		fmt.Printf("%s\n", out)
	}
	os.Exit(0)
}
