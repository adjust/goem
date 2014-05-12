package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var cmdTest = &Command{
	Run:  test,
	Name: "test",
}

func test(args []string) {
	testDir := config.Testdir
	setGoPath()
	buildPackages := false

	if testDir == "" {
		testDir = config.Testdir
	} else {
		dirs := strings.Split(testDir, " ")
		if len(dirs) > 1 && dirs[0] == "-i" {
			buildPackages = true
		}
		testDir = dirs[len(dirs)-1]
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
		fmt.Printf("%s %s\n", out, err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("%s\n", out)
	}
	os.Exit(0)
}
