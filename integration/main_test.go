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
  - name: input1
  - name: input2
outputs:
  - name: output1
  - name: output2

run:
  path: task.sh
`))
		Expect(err).ToNot(HaveOccurred())
		taskYaml.Close()
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()

		os.Remove(taskYaml.Name())
	})

	It("generates a shell script", func() {
		command := exec.Command(binaryPath, "-taskYamlPath", taskYaml.Name())
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out.Contents()).To(BeEquivalentTo(`#!/bin/bash

set -eu

TARGET=

INPUT1=$(mktemp -d -t input1)
INPUT2=$(mktemp -d -t input2)

OUTPUT1=$(mktemp -d -t output1)
OUTPUT2=$(mktemp -d -t output2)

fly \
  -t $TARGET \
  execute \
  -i input1=$INPUT1 \
  -i input2=$INPUT2 \
  -o output1=$OUTPUT1 \
  -o output2=$OUTPUT2 \
  -c task.yml

ls -l $OUTPUT1
ls -l $OUTPUT2

rm -rf $INPUT1
rm -rf $INPUT2

rm -rf $OUTPUT1
rm -rf $OUTPUT2
`))
	})

})
