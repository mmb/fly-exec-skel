package flyexecskel

import (
	"github.com/concourse/atc"
	"path"
	"strings"
)

func TaskInputName(taskConfig atc.TaskConfig) string {
	return strings.SplitN(taskConfig.Run.Path, "/", 2)[0]
}

func NonTaskInputs(taskConfig atc.TaskConfig) []atc.TaskInputConfig {
	taskInputName := TaskInputName(taskConfig)

	nonTaskInputs := make([]atc.TaskInputConfig, 0)
	for _, input := range taskConfig.Inputs {
		if input.Name != taskInputName {
			nonTaskInputs = append(nonTaskInputs, input)
		}
	}

	return nonTaskInputs
}

func RunPathToTaskInput(taskConfig atc.TaskConfig) string {
	depth := strings.Count(taskConfig.Run.Path, "/") - 1

	dirs := make([]string, 0)
	for i := 0; i < depth; i++ {
		dirs = append(dirs, "..")
	}

	return path.Join(dirs...)
}
