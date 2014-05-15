package main

import (
	"os"
)

var cmdBundle = &Command{
	Run:  bundle,
	Name: "bundle",
}

var bundled = map[string]string{}

func bundle(args []string) {
	makeBase()
	installDeps("")
	resolveDeps(args)
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

func getPackages(packages []*Package) {
	for _, pkg := range packages {
		if pkg.branchIsPath() {
			pkg.createSymlink()
			continue
		}
		if bundled[pkg.Name] == pkg.Branch {
			continue
		}
		if pkg.sourceExist() {
			pkg.updateSource()
		} else {
			pkg.getSource()
		}
		pkg.setHead()
		bundled[pkg.Name] = pkg.Branch
	}
}
