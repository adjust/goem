package main

import (
	"fmt"
	"os"
	"path"
)

var cmdBundle = &Command{
	Run:  Bundle,
	Name: "bundle",
}

func Bundle(args []string) {
	bundle(args)
	dep := NewDepCheck(config, args[0])
	dep.Start()
	config = NewLockConfig()
	bundle(args)
}

func bundle(args []string) {
	err := makeBase()
	found := false
	if err != nil {
		fmt.Printf("the following error occured while bundling\n")
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
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

func makeBase() error {
	goDirs := [3]string{"/src", "/pkg", "/bin"}
	for _, ext := range goDirs {
		err := os.MkdirAll(getGoPath()+ext, 0777)
		if err != nil {
			return fmt.Errorf("while creating needed dirs: " + err.Error())
		}
	}
	return nil
}

// getPackages() checks if there is already the wanted source
// if so it updates the source and sets the head as specified in the Gofile
// otherwise it fetches the source and sets the head
// on error getPackages exits
func getPackages(packages []Package) {
	for _, pkg := range packages {
		if pkg.BranchIsPath() {
			name, err := os.Getwd()
			if err != nil {
				fmt.Printf(err.Error())
			}
			name += "/" + pkg.Branch
			err = os.RemoveAll(getGoPath() + "/src/" + pkg.Name)
			if err != nil {
				fmt.Printf("while trying to remove useless folder: %s", err.Error())
			}
			os.Mkdir(getGoPath()+"/src/", 0777)
			err = os.Symlink(name, getGoPath()+"/src/"+pkg.Name)
			if err != nil {
				fmt.Printf("while trying to set 'self' link: %s\n", err.Error())
			}
			continue
		}
		if pkg.Branch == "self" {
			name, err := os.Getwd()
			if err != nil {
				fmt.Printf(err.Error())
			}
			name += "/" + path.Base(pkg.Name)
			err = os.RemoveAll(getGoPath() + "/src/" + pkg.Name)
			if err != nil {
				fmt.Printf("while trying to remove useless folder: %s", err.Error())
			}
			os.Mkdir(getGoPath()+"/src/"+pkg.Name, 0777)
			err = os.Symlink(name, getGoPath()+"/src/"+pkg.Name+"/"+path.Base(pkg.Name))
			if err != nil {
				fmt.Printf("while trying to set 'self' link: %s\n", err.Error())
			}
			continue
		}
		if !pkg.SourceExist() {
			err := getSource(pkg)
			if err != nil {
				fmt.Printf("while trying to get the source files: %s\n\n", err.Error())
				os.Exit(1)
			}
		} else {
			err := updateSource(pkg)
			if err != nil {
				fmt.Printf("while trying to update the source files: %s\n\n", err.Error())
				os.Exit(1)
			}
		}
		setHead(pkg)
	}
}

func getSource(pkg Package) error {
	err := git.clone(pkg)
	if err != nil {
		delErr := os.RemoveAll(path.Dir(getGoPath() + "/src/" + pkg.Name))
		if delErr != nil {
			fmt.Printf("An error occured while getting the source, but i am unable to tidy up\n")
			fmt.Printf("Please remove %s manually\n\n", path.Dir(getGoPath()+"/src/"+pkg.Name))
		}
		return fmt.Errorf("Could not get source:\n\n" + err.Error())
	}
	return nil
}

// updateSource() tries to call git pull after switching to master branch
// on error it returns the error
func updateSource(pkg Package) error {
	err := git.pull(pkg)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	return nil
}

// setHead() tries to set the head according to the Gofile with git checkout
// on error it returns the error
func setHead(pkg Package) error {
	err := git.checkout(pkg, "")
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	return nil
}
