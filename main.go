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
{{ with .TaskConfig -}}
{{ range $k, $v := .Params }}
{{ if $v -}}
# export {{ $k }}={{ $v -}}
{{ else -}}
# TODO set {{ $k }}
# export {{ $k }}=
echo ${{ $k -}}
{{ end -}}
{{ end -}}
{{ end }}

{{ divider "inputs" }}
{{ with .TaskConfig -}}
{{ range nonTaskInputs . }}
{{ envVarName .Name }}=$(mktemp -d -t {{ .Name }})
# TODO create test input in ${{ envVarName .Name -}}
{{ end }}

{{ divider "outputs" }}
{{ range .Outputs }}
{{ envVarName .Name }}=$(mktemp -d -t {{ .Name }})
{{- end -}}
{{ end }}

{{ divider "execute" }}

fly \
  -t {{ .Target }} \
  execute \
  -i {{ taskInputName .TaskConfig }}={{ runPathToTaskInput .TaskConfig }} \
{{ with .TaskConfig -}}
{{ range nonTaskInputs . }}  -i {{ .Name }}=${{ envVarName .Name }} \
{{ end -}}
{{ range .Outputs }}  -o {{ .Name }}=${{ envVarName .Name }} \
{{ end }}  -c task.yml

{{ divider "show outputs" }}
{{ range .Outputs }}
ls -l ${{ envVarName .Name -}}
{{ end }}

{{ divider "cleanup" }}
{{ range nonTaskInputs . }}
rm -rf ${{ envVarName .Name -}}
{{ end -}}
{{ range .Outputs }}
rm -rf ${{ envVarName .Name -}}
{{ end }}
{{ end }}`
	tmpl := template.New("script")
	tmpl.Funcs(template.FuncMap{
		"divider":            flyexecskel.Divider,
		"envVarName":         flyexecskel.EnvVarName,
		"nonTaskInputs":      flyexecskel.NonTaskInputs,
		"runPathToTaskInput": flyexecskel.RunPathToTaskInput,
		"taskInputName":      flyexecskel.TaskInputName,
	})
	tmpl.Parse(templateText)

	err = tmpl.Execute(os.Stdout, ti)
	if err != nil {
		panic(err)
	}
}
