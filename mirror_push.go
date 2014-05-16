package main

import (
	"fmt"
	"strings"
)

var cmdMirrorPush = &Command{
	Run:  mirrorPush,
	Name: "mirror-push",
}

func mirrorPush(args []string) {
	config.parse("")
	repos := dirRead(0, srcPath, nil)
	for _, repo := range repos {
		remote := fmt.Sprintf("%s.git", strings.Replace(repo, "github.com/", "", 1))
		fmt.Printf("pushing %s\n", repo)
		git.push(repo, config.Mirror, remote)
	}
}
