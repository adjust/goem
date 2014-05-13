package main

import (
	"os"
	"os/exec"
)

var cmdBuild = &Command{
	Run:  build,
	Name: "build",
}

func build(args []string) {
	binName := args[0]
	setGoPath()
	if binName == "" {
		binName = "a.out"
	}

	binName = strings.TrimSpace(binName)
	binName = strings.Replace(binName, "\n", "", -1)

	sourceFiles, err := getSourceFiles()
	if err != nil {
		fmt.Printf("while trying to collect source files: " + err.Error())
	}
	myArgs := []string{}
	myArgs = append(myArgs, "build")
	myArgs = append(myArgs, "-o")
	myArgs = append(myArgs, binName)
	myArgs = append(myArgs, sourceFiles...)

	execBuild := exec.Command(
		"go",
		myArgs...,
	)
	out, err := execBuild.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", out)
	}
}

func getSourceFiles() ([]string, error) {
	cwd, err := os.Getwd()
	var glob []string
	if err != nil {
		return nil, fmt.Errorf("while trying to get working dir: " + err.Error())
	}
	sourceFiles, err := filepath.Glob(cwd + "/*\\.go")
	if err != nil {
		return nil, fmt.Errorf("while trying to get glob filepath: " + err.Error())
	}
	regex := regexp.MustCompile("^\\.go")
	for _, file := range sourceFiles {
		base := path.Base(file)
		if !regex.MatchString(base) {
			glob = append(glob, base)
		}
	}
	return glob, nil
}
