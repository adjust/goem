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
	repos, err := dirRead(0, srcPath, nil)
	if err != nil {
		stderrAndExit(err)
	}
	cleaned := []string{}
	for i, repo := range repos {
		if repo == "" {
			continue
		}
		repos[i] = strings.Replace(repos[i], "github.com/", "", 1)
		cleaned = append(cleaned, repos[i])
	}
	templateGitoliteConfig(cleaned)
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
