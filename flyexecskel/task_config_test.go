package flyexecskel_test

import (
	"github.com/concourse/atc"
	"github.com/mmb/fly-exec-skel/flyexecskel"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TaskConfig", func() {
	Describe("TaskInputName", func() {
		Context("when the run is a relative path", func() {
			It("returns the name of the task input", func() {
				taskConfig := atc.TaskConfig{
					Run: atc.TaskRunConfig{Path: "task-repo/a/b/task.sh"},
				}
				Expect(flyexecskel.TaskInputName(taskConfig)).To(Equal("task-repo"))
			})
			Context("when the run is an absolute path", func() {
				It("returns the empty string", func() {
					taskConfig := atc.TaskConfig{
						Run: atc.TaskRunConfig{Path: "/a/b/task.sh"},
					}
					Expect(flyexecskel.TaskInputName(taskConfig)).To(Equal(""))
				})
			})
			Context("when the run has no directories", func() {
				It("returns the empty string", func() {
					taskConfig := atc.TaskConfig{
						Run: atc.TaskRunConfig{Path: "bash"},
					}
					Expect(flyexecskel.TaskInputName(taskConfig)).To(Equal(""))
				})
			})
		})
	})

	Describe("NonTaskInputs", func() {
		It("returns all inputs that are not the task input", func() {
			taskConfig := atc.TaskConfig{
				Run: atc.TaskRunConfig{Path: "task-repo/a/b/task.sh"},
				Inputs: []atc.TaskInputConfig{
					atc.TaskInputConfig{Name: "input1"},
					atc.TaskInputConfig{Name: "input2"},
					atc.TaskInputConfig{Name: "task-repo"},
				},
			}
			Expect(flyexecskel.NonTaskInputs(taskConfig)).To(Equal(
				[]atc.TaskInputConfig{
					atc.TaskInputConfig{Name: "input1"},
					atc.TaskInputConfig{Name: "input2"},
				},
			))
		})
	})

	Describe("RunPathToTaskInput", func() {
		It("returns the relative path from the run path to the task input root", func() {
			taskConfig := atc.TaskConfig{
				Run: atc.TaskRunConfig{Path: "task-repo/a/b/task.sh"},
			}
			Expect(flyexecskel.RunPathToTaskInput(taskConfig)).To(Equal("../.."))
		})
	})
})
