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

{{ divider "params" }}

{{ range $k, $v := .Params -}}
{{ if $v -}}
# export {{ $k }}={{ $v }}
{{ else -}}
# TODO set {{ $k }}
# export {{ $k }}=
echo ${{ $k }}
{{ end -}}
{{ end -}}
{{ end -}}

{{ if nonTaskInputs . -}}
{{ "\n" }}{{ divider "inputs" }}

{{ range nonTaskInputs . -}}
{{ envVarName .Name }}=$(mktemp -d -t {{ .Name }})
# TODO create test input in ${{ envVarName .Name }}
{{ end -}}
{{ end -}}

{{ if .Outputs -}}
{{ "\n" }}{{ divider "outputs" }}

{{ range .Outputs -}}
{{ envVarName .Name }}=$(mktemp -d -t {{ .Name }})
{{ end -}}
{{ end -}}

{{ end -}}

{{ "\n" }}{{ divider "execute" }}

fly \
  -t {{ .Target }} \
  execute \
  -i {{ taskInputName .TaskConfig }}={{ runPathToTaskInput .TaskConfig }} \
{{ with .TaskConfig -}}
{{ range nonTaskInputs . }}  -i {{ .Name }}=${{ envVarName .Name }} \
{{ end -}}
{{ range .Outputs }}  -o {{ .Name }}=${{ envVarName .Name }} \
{{ end }}  -c task.yml

{{ if .Outputs -}}
{{ divider "show outputs" }}

{{ range .Outputs -}}
ls -l ${{ envVarName .Name }}
{{ end -}}
{{ end -}}

{{ if (nonTaskInputs .) or .Outputs -}}
{{ "\n" }}{{ divider "cleanup" }}
{{ range nonTaskInputs . }}
rm -rf ${{ envVarName .Name -}}
{{ end -}}
{{ range .Outputs }}
rm -rf ${{ envVarName .Name -}}
{{ end -}}
{{ end -}}
{{ end }}
`
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
