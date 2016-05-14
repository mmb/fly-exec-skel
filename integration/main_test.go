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
	var taskYaml *os.File

	BeforeSuite(func() {
		var err error
		binaryPath, err = gexec.Build("github.com/mmb/fly-exec-skel")
		Expect(err).ToNot(HaveOccurred())

		taskYaml, err = ioutil.TempFile("", "task.yml")
		Expect(err).ToNot(HaveOccurred())

		_, err = taskYaml.Write([]byte(`---
platform: linux
inputs:
  - name: input-1
  - name: input-2
outputs:
  - name: output-1
  - name: output-2
run:
  path: task.sh
params:
  PARAM_1: param-1-default
  PARAM_2: param-2-default
  PARAM_3:
  PARAM_4:
`))
		Expect(err).ToNot(HaveOccurred())
		taskYaml.Close()
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()

		os.Remove(taskYaml.Name())
	})

	It("generates a shell script", func() {
		command := exec.Command(binaryPath, "-taskYamlPath", taskYaml.Name(), "-target", "test-target")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out.Contents()).To(BeEquivalentTo(`#!/bin/bash

set -eu

# export PARAM_1=param-1-default

# export PARAM_2=param-2-default

# export PARAM_3=<set PARAM_3 value>
echo $PARAM_3

# export PARAM_4=<set PARAM_4 value>
echo $PARAM_4

INPUT_1=$(mktemp -d -t input-1)
# Create test input in $INPUT_1

INPUT_2=$(mktemp -d -t input-2)
# Create test input in $INPUT_2

OUTPUT_1=$(mktemp -d -t output-1)

OUTPUT_2=$(mktemp -d -t output-2)

fly \
  -t test-target \
  execute \
  -i input-1=$INPUT_1 \
  -i input-2=$INPUT_2 \
  -o output-1=$OUTPUT_1 \
  -o output-2=$OUTPUT_2 \
  -c task.yml

ls -l $OUTPUT_1

ls -l $OUTPUT_2

rm -rf $INPUT_1

rm -rf $INPUT_2

rm -rf $OUTPUT_1

rm -rf $OUTPUT_2
`))
	})

})
