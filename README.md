[![Build Status](https://travis-ci.org/mmb/fly-exec-skel.svg?branch=master)](https://travis-ci.org/mmb/fly-exec-skel)

# fly-exec-skel
Generates skeleton Concourse fly execute shell scripts from a task YAML. These
scripts provide a reproducible workflow for fly executing a task with test
params, inputs and outputs.

Script workflow:

- set params, default params will be commented out, fill in your own values for
  required params

- create temporary input directories, fill in own code to stub the inputs
you need (create files, git init, etc.)

- create temporary output directories

- fly execute with temporary inputs and outputs

- show / test the outputs fly put in the output directories

- clean up up temporary directories

Scripts are expected to be run from the directory that the `task.yml` is in.

Example generated script:

```sh
#!/bin/bash

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
```
