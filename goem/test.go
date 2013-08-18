package goem

import (
	"fmt"
	"os"
	"os/exec"
)

func test(config *Config, testDir string) bool {
	setGoPath()

	if testDir == "" {
		testDir = config.Testdir
	}
	os.Chdir(testDir)

	cmd := exec.Command(
		"go",
		"test",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s %s\n", out, err.Error())
	}
	cmd = exec.Command(
		"go",
		"test",
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s %s\n", out, err.Error())
		return false
	} else {
		fmt.Printf("%s\n", out)
	}
	return true
}
