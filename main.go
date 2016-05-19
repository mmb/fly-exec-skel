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

{{ divider "params" }}
{{ with .TaskConfig }}{{ range $k, $v := .Params }}
{{ if $v }}# export {{ $k }}={{ $v }}
{{ else }}# uncomment and set {{ $k }} value
# export {{ $k }}=
echo ${{ $k }}
{{ end }}{{ end }}{{ end }}{{ with .TaskConfig }}
{{ divider "inputs" }}
{{ range .Inputs }}
{{ envVarName .Name }}=$(mktemp -d -t {{ .Name }})
# Create test input in ${{ envVarName .Name }}
{{ end }}
{{ divider "outputs" }}
{{ range .Outputs }}
{{ envVarName .Name }}=$(mktemp -d -t {{ .Name }})
{{ end }}{{ end }}
{{ divider "execute" }}

fly \
  -t {{ .Target }} \
  execute \
{{ with .TaskConfig }}{{ range .Inputs }}  -i {{ .Name }}=${{ envVarName .Name }} \
{{ end }}{{ range .Outputs }}  -o {{ .Name }}=${{ envVarName .Name }} \
{{ end }}  -c task.yml

{{ divider "show outputs" }}
{{ range .Outputs }}
ls -l ${{ envVarName .Name }}
{{ end }}
{{ divider "cleanup" }}
{{ range .Inputs }}
rm -rf ${{ envVarName .Name }}
{{ end }}{{ range .Outputs }}
rm -rf ${{ envVarName .Name }}
{{ end }}{{ end }}`
	tmpl := template.New("script")
	tmpl.Funcs(template.FuncMap{
		"divider":    flyexecskel.Divider,
		"envVarName": flyexecskel.EnvVarName,
	})
	tmpl.Parse(templateText)

	err = tmpl.Execute(os.Stdout, ti)
	if err != nil {
		panic(err)
	}
}
