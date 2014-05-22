package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var cmdMirrorConfig = &Command{
	Run:  mirrorConfig,
	Name: "mirror-config",
}

const tmpl = `
{{range $index, $elem := .}}
repo    {{ $elem }}
        RW+     =   @all
{{end}}
`

func mirrorConfig(args []string) {
	repos := dirRead(0, srcPath, nil)
	for i, _ := range repos {
		repos[i] = strings.Replace(repos[i], "github.com/", "", 1)
	}
	templateGitoliteConfig(repos)
}

func templateGitoliteConfig(repos []string) {
	var buf bytes.Buffer
	configTemplate, err := template.New("gitolite").Parse(tmpl)
	if err != nil {
		stderrAndExit(err)
	}
	err = configTemplate.Execute(&buf, repos)
	fmt.Println(buf.String())
}
