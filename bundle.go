package main

import (
	"os"
)

var cmdBundle = &Command{
	Run:  Bundle,
	Name: "bundle",
}

func Bundle(args []string) {
	makeBase()
	installDeps("")
	ResolveDeps(args)
	installDeps("Gofile.lock")
}

func makeBase() {
	goDirs := [3]string{"/src", "/pkg", "/bin"}
	for _, ext := range goDirs {
		err := os.MkdirAll(getGoPath()+ext, 0777)
		if err != nil {
			stderrAndExit(err)
		}
	}
}

func installDeps(gofile string) {
	config.parse(gofile)
	found := false
	for _, env := range config.Env {
		if getGoEnv() == env.Name {
			getPackages(env.Packages)
			found = true
			break
		}
	}
	if !found {
		for _, env := range config.Env {
			if "development" == env.Name {
				getPackages(env.Packages)
				found = true
				break
			}
		}
	}
}

func getPackages(packages []Package) {
	for _, pkg := range packages {
		if pkg.BranchIsPath() {
			pkg.CreateSymlink()
			continue
		}
		if pkg.SourceExist() {
			pkg.UpdateSource()
		} else {
			pkg.GetSource()
		}
		pkg.SetHead()
	}
}
