package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var cmdBuild = &Command{
	Run:  build,
	Name: "build",
}

func build(args []string) {
	var binName = "a.out"
	if len(args) > 0 {
		binName = args[0]
	}
	setGoPath()

	binName = strings.TrimSpace(binName)
	binName = strings.Replace(binName, "\n", "", -1)

	cmdArgs := []string{"build", "-o", binName}
	cmdArgs = append(cmdArgs, getSourceFiles()...)
	execBuild := exec.Command(
		"go",
		cmdArgs...,
	)

	out, err := execBuild.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}
}

func getSourceFiles() []string {
	cwd, err := os.Getwd()
	if err != nil {
		stderrAndExit(err)
	}

	sourceFiles, err := filepath.Glob(cwd + `/*\.go$`)
	if err != nil {
		stderrAndExit(err)
	}

	regex := regexp.MustCompile(`^\.go$`)
	var glob []string
	for _, file := range sourceFiles {
		base := path.Base(file)
		if !regex.MatchString(base) {
			glob = append(glob, base)
		}
	}
	return glob
}
