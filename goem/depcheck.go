package goem

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

// DepCheck object is used to determine recursive dependencies
// it tries to combine all dependencies of the dependencies into
// Gofile lock for the given environment
type DepCheck struct {
	// the GOPATH environment var
	goPath string
	// a goem Lister Object to resolve the dependencies
	list *Lister
	// the config this package
	config *Config
}

// Constructor that sets the defaults
func NewDepCheck(config *Config) *DepCheck {
	return &DepCheck{
		goPath: getGoPath(),
		list:   NewList(),
		config: config,
	}
}

// Start() does the following steps:
//   1. get the gofiles of the dependencies
//   2. parse the gofiles and determine dependencies
//   3. write these dependencies to Gofile.lock
func (self *DepCheck) Start() {
	pkgMap := map[string]string{}
	for {
		before := len(pkgMap)
		subConfigs := self.getGofiles()
		self.checkDeps(subConfigs, pkgMap)
		after := len(pkgMap)
		if before == after {
			break
		}
	}
	self.writeGofileLock(pkgMap)
}

// getGoFiles() does the following steps
//   1. get all installed packages with the lister object
//   2. check if a Gofile for this package exists
//   3. return a list of installed packages with corresponding configs
// on error it exits with 1
func (self *DepCheck) getGofiles() []*GoPkg {
	packages, err := self.list.dirRead(0, self.goPath+"/src", nil)
	if err != nil {
		fmt.Printf("Something went wrong: %s", err.Error())
		os.Exit(1)
	}
	goPkgs := make([]*GoPkg, len(packages))
	for i, pkg := range packages {
		if pkg == "" {
			continue
		}
		if self.fileExists(self.goPath + "src/" + pkg + "/Gofile") {
			goPkg := &GoPkg{
				name:   pkg,
				config: &Config{},
			}
			goPkg.config.parse(self.goPath + "src/" + pkg + "/Gofile")
			goPkgs[i] = goPkg
		} else {
			fmt.Printf("Did not find a Gofile for: %s\n", pkg)
		}
	}
	return self.shrink(goPkgs)
}

// shrink() removes empty entries from an array of GoPkg
func (self *DepCheck) shrink(configs []*GoPkg) []*GoPkg {
	realLen := 0
	for _, config := range configs {
		if config != nil {
			realLen++
		}
	}
	newConfigs := make([]*GoPkg, realLen+1)
	index := 0
	for _, config := range configs {
		if config != nil {
			newConfigs[index] = config
			index++
		}
	}
	newConfigs[index] = &GoPkg{name: "this", config: self.config}
	return newConfigs
}

// fileExists() checks if the given file exists
func (self *DepCheck) fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// checkDeps calls checkDeep() for all Packages with a Gofile to find conflicts
// and all dependencies
func (self *DepCheck) checkDeps(goPkgs []*GoPkg, pkgMap map[string]string) map[string]string {
	for _, goPkg := range goPkgs {
		self.checkDeep(goPkg, goPkgs, pkgMap)
	}
	return pkgMap
}

// checkDeep() iterates over all installed Packages and compares them with
// the other installed packages
func (self *DepCheck) checkDeep(goPkg *GoPkg, goPkgs []*GoPkg, pkgMap map[string]string) {
	for _, otherGoPkg := range goPkgs {
		otherPkgs := self.getPkgForEnv(otherGoPkg)
		pkgs := self.getPkgForEnv(goPkg)
		self.cmpPkgList(pkgs, otherPkgs, goPkg, otherGoPkg, pkgMap)
	}
}

// getPkgForEnv returnes a list of all packages for the given environment
func (self *DepCheck) getPkgForEnv(goPkg *GoPkg) []Package {
	if len(goPkg.config.Env) == 1 {
		goPkg.config.Env[0].Name = getGoEnv()
	}
	for _, env := range goPkg.config.Env {
		if env.Name == getGoEnv() {
			return env.Packages
		}
	}
	return nil
}

// cmpPkgList compares the dependencies of 2 packages
// it prints errors, which it cannot resolve to the console
// it changes the given pkgMap to save its results
func (self *DepCheck) cmpPkgList(pkgList, otherPkgList []Package, goPkg, otherGoPkg *GoPkg, pkgMap map[string]string) {
	sort.Sort(&ByName{pkgList})
	sort.Sort(&ByName{otherPkgList})
	counter1 := len(pkgList)
	counter2 := len(otherPkgList)
	iter1 := 0
	iter2 := 0
	for {
		if counter1 == 0 || counter2 == 0 {
			break
		}
		if pkgList[iter1].Name == otherPkgList[iter2].Name {
			if pkgList[iter1].Branch != otherPkgList[iter2].Branch {
				counter1--
				counter2--
				result := self.resolveDep(pkgList[iter1], otherPkgList[iter2])
				if result != "" {
					pkgMap[pkgList[iter1].Name] = result
				} else {
					fmt.Printf("Cannot resolve dependency:\n")
					fmt.Printf(
						"Package %s needed in Version: %s for %s and in %s for %s\n",
						pkgList[iter1].Name,
						otherPkgList[iter2].Branch,
						otherGoPkg.name,
						pkgList[iter1].Branch,
						goPkg.name,
					)
				}
				continue
			} else {
				pkgMap[pkgList[iter1].Name] = pkgList[iter1].Branch
				pkgMap[otherPkgList[iter2].Name] = otherPkgList[iter2].Branch
			}
		}
		if pkgList[iter1].Name < otherPkgList[iter2].Name {
			counter1--
			iter1++
		} else {
			counter2--
			iter2++
		}
	}
}

func (self *DepCheck) resolveDep(pkg1, pkg2 Package) string {
	branch1, _ := git.refNameToCommit(pkg1)
	if branch1 == "" {
		branch1 = pkg1.Branch
	}
	branch2, _ := git.refNameToCommit(pkg2)
	if branch2 == "" {
		branch2 = pkg2.Branch
	}
	if branch1 == branch2 {
		return branch1
	}

	okay := true
	if branch1[0] == '<' {
		//check if branch2 is < than branch1
		if branch2[0] == '>' || branch2[0] == '<' {
			if !self.isOlderThan(branch2[1:], branch1[1:], pkg1) {
				okay = false
			}
		} else {
			if !self.isOlderThan(branch2, branch1[1:], pkg1) {
				okay = false
			}
		}
	}
	if branch1[0] == '>' {
		//check if branch2 is > than branch1
		if branch2[0] == '>' || branch2[0] == '<' {
			if !self.isOlderThan(branch1[1:], branch2[1:], pkg1) {
				okay = false
			}
		} else {
			if !self.isOlderThan(branch1[1:], branch2, pkg1) {
				okay = false
			}
		}
	}
	if okay {
		return branch1
	}
	return ""
}

func (self *DepCheck) isOlderThan(branch1, branch2 string, pkg Package) bool {
	gitLog, err := git.log(pkg, "")
	if err != nil {
		fmt.Printf(err.Error())
	}
	age1 := 0
	age2 := 0
	for i, line := range gitLog {
		if line == branch1 {
			age1 = i
		}
		if line == branch2 {
			age2 = i
		}
	}
	return age1 > age2
}

// writeGofileLock writes the collected dependencies to Gofile.Lock
func (self *DepCheck) writeGofileLock(deps map[string]string) {
	packages := make([]Package, len(deps))
	iter := 0
	for name, branch := range deps {
		pkg := Package{Name: name, Branch: branch}
		packages[iter] = pkg
		iter++
	}
	env := Env{Name: getGoEnv(), Packages: packages}
	content := []Env{env}
	config := &Config{Env: content}
	j, _ := json.MarshalIndent(config, "    ", "    ")
	ioutil.WriteFile("./Gofile.lock", j, 0777)
}
