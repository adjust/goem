package goem

import (
	"fmt"
	"os"
	"os/exec"
)

func test(config *Config) bool {
	setGoPath()
	os.Chdir(config.Testdir)
	cmd := exec.Command(
		"go",
		"test",
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
