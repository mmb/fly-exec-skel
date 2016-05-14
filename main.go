package main

import (
	"flag"
	"github.com/concourse/atc"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type templateInput struct {
	TaskConfig atc.TaskConfig
	Target     string
}

func main() {
	var taskYamlPath string
	ti := new(templateInput)

	flag.StringVar(&taskYamlPath, "taskYamlPath", "./task.yml", "path to task YAML")
	flag.StringVar(&ti.Target, "target", "private", "Concourse target name")
	flag.Parse()

	taskYamlBytes, err := ioutil.ReadFile(taskYamlPath)
	if err != nil {
		panic(err)
	}
	ti.TaskConfig, err = atc.LoadTaskConfig(taskYamlBytes)
	if err != nil {
		panic(err)
	}

	templateText := `#!/bin/bash

set -eu
{{ with .TaskConfig }}{{ range .Inputs }}
{{ upcase .Name }}=$(mktemp -d -t {{ .Name }})
# Create test input in ${{ upcase .Name }}
{{ end }}{{ range .Outputs }}
{{ upcase .Name }}=$(mktemp -d -t {{ .Name }})
{{ end }}{{ end }}
fly \
  -t {{ .Target }} \
  execute \
{{ with .TaskConfig }}{{ range .Inputs }}  -i {{ .Name }}=${{ upcase .Name }} \
{{ end }}{{ range .Outputs }}  -o {{ .Name }}=${{ upcase .Name }} \
{{ end }}  -c task.yml
{{ range .Outputs }}
ls -l ${{ upcase .Name }}
{{ end }}{{ range .Inputs }}
rm -rf ${{ upcase .Name }}
{{ end }}{{ range .Outputs }}
rm -rf ${{ upcase .Name }}
{{ end }}{{ end }}`
	tmpl := template.New("script")
	tmpl.Funcs(template.FuncMap{"upcase": strings.ToUpper})
	tmpl.Parse(templateText)

	err = tmpl.Execute(os.Stdout, ti)
	if err != nil {
		panic(err)
	}
}
