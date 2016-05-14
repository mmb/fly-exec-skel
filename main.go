package main

import (
	"flag"
	"github.com/concourse/atc"
	"github.com/mmb/fly-exec-skel/flyexecskel"
	"io/ioutil"
	"os"
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
{{ envVarName .Name }}=$(mktemp -d -t {{ .Name }})
# Create test input in ${{ envVarName .Name }}
{{ end }}{{ range .Outputs }}
{{ envVarName .Name }}=$(mktemp -d -t {{ .Name }})
{{ end }}{{ end }}
fly \
  -t {{ .Target }} \
  execute \
{{ with .TaskConfig }}{{ range .Inputs }}  -i {{ .Name }}=${{ envVarName .Name }} \
{{ end }}{{ range .Outputs }}  -o {{ .Name }}=${{ envVarName .Name }} \
{{ end }}  -c task.yml
{{ range .Outputs }}
ls -l ${{ envVarName .Name }}
{{ end }}{{ range .Inputs }}
rm -rf ${{ envVarName .Name }}
{{ end }}{{ range .Outputs }}
rm -rf ${{ envVarName .Name }}
{{ end }}{{ end }}`
	tmpl := template.New("script")
	tmpl.Funcs(template.FuncMap{"envVarName": flyexecskel.EnvVarName})
	tmpl.Parse(templateText)

	err = tmpl.Execute(os.Stdout, ti)
	if err != nil {
		panic(err)
	}
}
