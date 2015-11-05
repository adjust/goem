package main

import (
	"fmt"
	"path"
	"sort"
)

var quiet bool

type GoPkg struct {
	name   string
	config *Config
}

func resolveDeps(packages Packages, args []string, mirrored bool) {
	if len(args) > 0 && args[0][0] == 'q' {
		quiet = true
	}

	pkgMap := map[string]string{}
	for {
		before := len(pkgMap)
		subConfigs := getGofiles(packages)
		checkDeps(subConfigs, pkgMap)
		after := len(pkgMap)
		if before == after {
			break
		}
		writeGofileLock(pkgMap)
		installDeps("Gofile.lock", mirrored)
	}
}

func getGofiles(packages Packages) []*GoPkg {
	goPath := getGoPath()
	goPkgs := make([]*GoPkg, len(packages))

	for i, pkg := range packages {
		gofilePath := path.Join(goPath, "src", pkg.Name, "Gofile")
		if fileExists(gofilePath) {
			goPkg := &GoPkg{
				name:   pkg.Name,
				config: &Config{},
			}
			goPkg.config.parse(gofilePath)
			goPkgs[i] = goPkg
		} else if !quiet {
			fmt.Printf("Did not find a Gofile for: %s\n", pkg.Name)
		}
	}

	return shrink(goPkgs)
}

func shrink(configs []*GoPkg) []*GoPkg {
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
	newConfigs[index] = &GoPkg{name: "this", config: config}
	return newConfigs
}

func checkDeps(goPkgs []*GoPkg, pkgMap map[string]string) map[string]string {
	for _, goPkg := range goPkgs {
		checkDeep(goPkg, goPkgs, pkgMap)
	}
	return pkgMap
}

func checkDeep(goPkg *GoPkg, goPkgs []*GoPkg, pkgMap map[string]string) {
	for _, otherGoPkg := range goPkgs {
		otherPkgs := getPkgForEnv(otherGoPkg)
		pkgs := getPkgForEnv(goPkg)
		cmpPkgList(pkgs, otherPkgs, goPkg, otherGoPkg, pkgMap)
	}
}

func getPkgForEnv(goPkg *GoPkg) []*Package {
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

func cmpPkgList(pkgList, otherPkgList []*Package, goPkg, otherGoPkg *GoPkg, pkgMap map[string]string) {
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
				result := resolveDep(pkgList[iter1], otherPkgList[iter2])
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

func resolveDep(pkg1, pkg2 *Package) string {
	if pkg1.branchIsPath() {
		return pkg1.Branch
	}
	if pkg2.branchIsPath() {
		return pkg2.Branch
	}
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
			if !isOlderThan(branch2[1:], branch1[1:], pkg1) {
				okay = false
			}
		} else {
			if !isOlderThan(branch2, branch1[1:], pkg1) {
				okay = false
			}
		}
	}
	if branch1[0] == '>' {
		//check if branch2 is > than branch1
		if branch2[0] == '>' || branch2[0] == '<' {
			if !isOlderThan(branch1[1:], branch2[1:], pkg1) {
				okay = false
			}
		} else {
			if !isOlderThan(branch1[1:], branch2, pkg1) {
				okay = false
			}
		}
	}
	if okay {
		return branch1
	}
	return ""
}

func isOlderThan(branch1, branch2 string, pkg *Package) bool {
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

func writeGofileLock(deps map[string]string) {
	packages := make([]*Package, len(deps))
	iter := 0
	for name, branch := range deps {
		pkg := &Package{Name: name, Branch: branch}
		packages[iter] = pkg
		iter++
	}
	env := &Env{Name: getGoEnv(), Packages: packages}
	content := []*Env{env}
	config := &Config{Env: content, Mirror: config.Mirror}
	config.write("./Gofile.lock")
}
