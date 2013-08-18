package goem

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func test(config *Config, testDir string) bool {
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

	cmd := exec.Command(
		"go",
		"test",
		testCommand,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s %s\n", out, err.Error())
		return false
	} else {
		fmt.Printf("%s\n", out)
	}
	return true
}
