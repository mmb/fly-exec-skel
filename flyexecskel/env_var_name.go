package flyexecskel

import (
	"fmt"
	"strings"
)

func InputEnvVarName(s string) string {
	return fmt.Sprintf("%s_INPUT", envVarName(s))
}

func OutputEnvVarName(s string) string {
	return fmt.Sprintf("%s_OUTPUT", envVarName(s))
}

func envVarName(s string) string {
	return strings.Replace(strings.ToUpper(s), "-", "_", -1)
}
