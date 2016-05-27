package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
)

var _ = Describe("Integration", func() {
	var binaryPath string
	var taskYamlFile *os.File
	taskYaml := `---
platform: linux
inputs:
  - name: input-1
  - name: input-2
  - name: task-repo
outputs:
  - name: output-1
  - name: output-2
run:
  path: task-repo/task1/task.sh
params:
  PARAM_1: param-1-default
  PARAM_2: param-2-default
  PARAM_3:
  PARAM_4:
`

	BeforeSuite(func() {
		var err error
		binaryPath, err = gexec.Build("github.com/mmb/fly-exec-skel")
		Expect(err).ToNot(HaveOccurred())

		taskYamlFile, err = ioutil.TempFile("", "task.yml")
		Expect(err).ToNot(HaveOccurred())

		_, err = taskYamlFile.Write([]byte(taskYaml))
		Expect(err).ToNot(HaveOccurred())
		taskYamlFile.Close()
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()

		os.Remove(taskYamlFile.Name())
	})

	It("generates a shell script", func() {
		command := exec.Command(binaryPath, "-taskYamlPath", taskYamlFile.Name(), "-target", "test-target")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out.Contents()).To(BeEquivalentTo(`#!/bin/bash

set -eu

# params -----------------------------------------------------------------------

# export PARAM_1=param-1-default
# export PARAM_2=param-2-default
# TODO set PARAM_3
# export PARAM_3=
echo $PARAM_3
# TODO set PARAM_4
# export PARAM_4=
echo $PARAM_4

# inputs -----------------------------------------------------------------------

INPUT_1=$(mktemp -d -t input-1)
# TODO create test input in $INPUT_1
INPUT_2=$(mktemp -d -t input-2)
# TODO create test input in $INPUT_2

# outputs ----------------------------------------------------------------------

OUTPUT_1=$(mktemp -d -t output-1)
OUTPUT_2=$(mktemp -d -t output-2)

# execute ----------------------------------------------------------------------

fly \
  -t test-target \
  execute \
  -i task-repo=.. \
  -i input-1=$INPUT_1 \
  -i input-2=$INPUT_2 \
  -o output-1=$OUTPUT_1 \
  -o output-2=$OUTPUT_2 \
  -c task.yml

# show outputs -----------------------------------------------------------------

ls -l $OUTPUT_1
ls -l $OUTPUT_2

# cleanup ----------------------------------------------------------------------

rm -rf $INPUT_1
rm -rf $INPUT_2
rm -rf $OUTPUT_1
rm -rf $OUTPUT_2
`))
	})

})
