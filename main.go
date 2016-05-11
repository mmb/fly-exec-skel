package main

import (
	"flag"
	"github.com/concourse/atc"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func main() {
	var taskYamlPath string

	flag.StringVar(&taskYamlPath, "taskYamlPath", "./task.yml", "path to task YAML")
	flag.Parse()

	taskYamlBytes, err := ioutil.ReadFile(taskYamlPath)
	if err != nil {
		panic(err)
	}
	taskConfig, err := atc.LoadTaskConfig(taskYamlBytes)
	if err != nil {
		panic(err)
	}

	templateText := `#!/bin/bash

set -eu

TARGET=
{{ range .Inputs }}
{{ upcase .Name }}=$(mktemp -d -t {{ .Name }}){{ end }}
{{ range .Outputs }}
{{ upcase .Name }}=$(mktemp -d -t {{ .Name }}){{ end }}

fly \
  -t $TARGET \
  execute \
{{ range .Inputs }}  -i {{ .Name }}=${{ upcase .Name }} \
{{ end }}{{ range .Outputs }}  -o {{ .Name }}=${{ upcase .Name }} \
{{ end }}  -c task.yml
{{ range .Outputs }}
ls -l ${{ upcase .Name }}{{ end }}
{{ range .Inputs }}
rm -rf ${{ upcase .Name }}{{ end }}
{{ range .Outputs }}
rm -rf ${{ upcase .Name }}{{ end }}
`
	tmpl := template.New("script")
	tmpl.Funcs(template.FuncMap{"upcase": strings.ToUpper})
	tmpl.Parse(templateText)

	err = tmpl.Execute(os.Stdout, taskConfig)
	if err != nil {
		panic(err)
	}
}
