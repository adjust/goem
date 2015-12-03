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
	installedPackages := installDeps("", false)
	resolveDeps(installedPackages, args, false)
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

func installDeps(gofile string, mirrored bool) Packages {
	config.parse(gofile)
	if mirrored {
		config.mirrored()
	}

	requestedEnvName := getGoEnv()

	var developmentEnv *Env
	for _, env := range config.Env {
		if env.Name == "development" {
			developmentEnv = env
		}

		if env.Name == requestedEnvName {
			return getPackages(env.Packages)
		}
	}

	if developmentEnv != nil {
		return getPackages(developmentEnv.Packages)
	}

	return nil
}

func getPackages(packages Packages) Packages {
	setGoPath()

	for _, pkg := range packages {
		if pkg.branchIsPath() {
			pkg.createSymlink()
			continue
		}
		if pkg.Branch != "" && bundled[pkg.Name] == pkg.Branch {
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

	return packages
}
