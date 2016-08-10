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
{{ with .TaskConfig -}}

{{ if .Params -}}

{{ "\n" }}{{ divider "params" }}

{{ range $k, $v := .Params -}}
{{ if $v -}}
# export {{ paramEnvVarName $k }}={{ $v }}
{{ else -}}
# TODO set {{ paramEnvVarName $k }}
# export {{ paramEnvVarName $k }}=
echo ${{ paramEnvVarName $k }}
{{ end -}}
{{ end -}}
{{ end -}}

{{ if nonTaskInputs . -}}
{{ "\n" }}{{ divider "inputs" }}

{{ range nonTaskInputs . -}}
{{ inputEnvVarName .Name }}=$(mktemp -d -t {{ .Name }})
# TODO create test input in ${{ inputEnvVarName .Name }}
{{ end -}}
{{ end -}}

{{ if .Outputs -}}
{{ "\n" }}{{ divider "outputs" }}

{{ range .Outputs -}}
{{ outputEnvVarName .Name }}=$(mktemp -d -t {{ .Name }})
{{ end -}}
{{ end -}}

{{ end -}}

{{ "\n" }}{{ divider "execute" }}

fly \
  -t {{ .Target }} \
  execute \
  -i {{ taskInputName .TaskConfig }}={{ runPathToTaskInput .TaskConfig }} \
{{ with .TaskConfig -}}
{{ range nonTaskInputs . }}  -i {{ .Name }}=${{ inputEnvVarName .Name }} \
{{ end -}}
{{ range .Outputs }}  -o {{ .Name }}=${{ outputEnvVarName .Name }} \
{{ end }}  -c task.yml{{ "\n" }}

{{- if .Outputs -}}
{{ "\n" }}{{ divider "show outputs" }}

{{ range .Outputs -}}
ls -l ${{ outputEnvVarName .Name }}
{{ end -}}
{{ end -}}

{{ if (nonTaskInputs .) or .Outputs -}}
{{ "\n" }}{{ divider "cleanup" }}

{{ range nonTaskInputs . -}}
rm -rf ${{ inputEnvVarName .Name }}
{{ end -}}
{{ range .Outputs -}}
rm -rf ${{ outputEnvVarName .Name }}
{{ end -}}
{{ end -}}
{{ end -}}
`
	tmpl := template.New("script")
	tmpl.Funcs(template.FuncMap{
		"divider":            flyexecskel.Divider,
		"inputEnvVarName":    flyexecskel.InputEnvVarName,
		"nonTaskInputs":      flyexecskel.NonTaskInputs,
		"outputEnvVarName":   flyexecskel.OutputEnvVarName,
		"paramEnvVarName":    flyexecskel.ParamEnvVarName,
		"runPathToTaskInput": flyexecskel.RunPathToTaskInput,
		"taskInputName":      flyexecskel.TaskInputName,
	})
	tmpl.Parse(templateText)

	err = tmpl.Execute(os.Stdout, ti)
	if err != nil {
		panic(err)
	}
}
