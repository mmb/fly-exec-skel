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
	TaskConfig   atc.TaskConfig
	Target       string
	TaskYamlPath string
}

func main() {
	ti := new(templateInput)

	flag.StringVar(&ti.TaskYamlPath, "taskYamlPath", "./task.yml", "path to task YAML")
	flag.StringVar(&ti.Target, "target", "private", "Concourse target name")
	flag.Parse()

	taskYamlBytes, err := ioutil.ReadFile(ti.TaskYamlPath)
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
{{ if taskInputName .TaskConfig -}}
{{ "  " }}-i {{ taskInputName .TaskConfig }}={{ runPathToTaskInput .TaskConfig }} \{{ "\n" }}
{{- end -}}

{{ range nonTaskInputs .TaskConfig }}  -i {{ .Name }}=${{ inputEnvVarName .Name }} \
{{ end -}}
{{ range .TaskConfig.Outputs }}  -o {{ .Name }}=${{ outputEnvVarName .Name }} \
{{ end }}  -c {{ .TaskYamlPath }}{{ "\n" }}

{{- if .TaskConfig.Outputs -}}
{{ "\n" }}{{ divider "show outputs" }}

{{ range .TaskConfig.Outputs -}}
ls -l ${{ outputEnvVarName .Name }}
{{ end -}}
{{ end -}}

{{ if (nonTaskInputs .TaskConfig) or .TaskConfig.Outputs -}}
{{ "\n" }}{{ divider "cleanup" }}

{{ range nonTaskInputs .TaskConfig -}}
rm -rf ${{ inputEnvVarName .Name }}
{{ end -}}
{{ range .TaskConfig.Outputs -}}
rm -rf ${{ outputEnvVarName .Name }}
{{ end -}}
{{ end -}}
`
	tmpl := template.New("script")
	tmpl.Funcs(template.FuncMap{
		"divider":            flyexecskel.Divider,
		"inputEnvVarName":    flyexecskel.InputEnvVarName,
		"nonTaskInputs":      flyexecskel.NonTaskInputs,
		"outputEnvVarName":   flyexecskel.OutputEnvVarName,
		"runPathToTaskInput": flyexecskel.RunPathToTaskInput,
		"taskInputName":      flyexecskel.TaskInputName,
	})
	tmpl.Parse(templateText)

	err = tmpl.Execute(os.Stdout, ti)
	if err != nil {
		panic(err)
	}
}
