package main

import (
	"fmt"
	"strings"
)

var cmdHelp = &Command{
	Usage: "help blah",
	Long:  "blah keys",
	Name:  "help",
}

func init() {
	cmdHelp.Run = Help
}

func Help(args []string) {
	if len(args) == 0 {
		helpHelp()
		return
	}
	switch strings.TrimSpace(args[0]) {
	case "init":
		helpInit()
	case "list":
		helpList()
	case "bundle":
		helpBundle()
	case "build":
		helpBuild()
	case "test":
		helpTest()
	default:
		helpHelp()
	}
}

func helpInit() {
	fmt.Printf(
		`
        usage: goem init

        Creates and initializes a new Gofile in current directory.

        This command will not override already existing Gofile.

`,
	)
}

func helpList() {
	fmt.Printf(
		`
        usage: goem list

        List lists all packages currently installed in your local gopath.

`,
	)
}

func helpBundle() {
	fmt.Printf(
		`
        usage: goem bundle

        Bundle reads your Gofile and fetches all packages in the desired version.

        An example Gofile looks like that:

        {
            "env" : [
                {
                    "name" : "development",
                    "packages": [
                        {
                            "name" : "github.com/adjust/goenv",
                            "branch" : "<cd3a33acbec38335c00b4ae252274827893d4e5b"
                        }
                    ]
                }
            ]
        }

        env:        This is your GO_ENV which is read from GO_ENV environment variable.
                    If the environment variable is not set, goem sets it to 'development'.
                    If the Gofile provides only one env, this env is read.

        packages:   Packages have a name and a branch. The name is the same as the one
                    you would use for a 'go get <package>. The branch can be a branch,
                    a commit hash or a tag.
                    If branch is set to 'self', goem assumes you have your library *.go
                    files in a subfolder with the same name as your project and symlinks this
                    to the goem Gopath.
             e.g.:
                For example the following directory structure:
                    project_name
                    |
                    |__project_name
                    |   |
                    |   |__
                    |
                    |__.go
                    |   |
                    |   |__src
                    |       |
                    |       |__your_account
                    |
                    |
                    |__main.go

               Would create a symlink from project_name/projectname to
               project_name/.go/src/your_account/project_name

`,
	)
}

func helpBuild() {
	fmt.Printf(
		`
        usage: goem build <binname>

        Build your project with all the dependencies listed in your Gofile.
        If no 'binname' is provided, the binary is called a.out.
        'binname' can also be a path:

        For example the following directory structure:
            project_root
            |
            |__bin
            |   |
            |   |__
            |
            |__main.go


        To build the binary in bin/ called 'nifty_bin', you would call
        'goem build bin/nifty_bin'

        Absolute paths are also allowed.

`,
	)
}

func helpTest() {
	fmt.Printf(
		`
        usage: goem test

        Run the tests as you would with go test.
        As goem commands need to be run from the projects root,
        you have to specify the test directory, if it is not the root directory.

        For example the following directory structure:
            project_root
            |
            |__test_folder
            |   |
            |   |_nifty_test.go
            |
            |__main.go

        In this case you need to add the 'testdir' key in your Gofile:
            e.g.:
                {
                    "testdir": "test_folder"
                }

        Then run goem test

`,
	)
}

func helpHelp() {
	fmt.Printf(
		`
        usage: goem help <topic>

        All goem commands only work on the root of your go project.
        The root is the directory in which your Gofile is located.
        Following commands are avaiable:
            - goem init
                Initialize a new Gofile
            - goem list
                List all bundled packages
            - goem build
                Build a binary
            - goem bundle
                Get all dependencies defined in your Gofile
            - goem test
                Run the tests
            - goem help
                Get additional help for the commands:
                e.g. goem help list

`,
	)
}
