package goem

import (
	"fmt"
	"os"
	"os/exec"
)

func test(config *Config) {
	setGoPath()
	os.Chdir(config.Testdir)
	cmd := exec.Command(
		"/usr/bin/go",
		"test",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s %s\n", out, err.Error())
	}
	cmd = exec.Command(
		"/usr/bin/go",
		"test",
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s %s\n", out, err.Error())
	} else {
		fmt.Printf("%s\n", out)
	}

}
