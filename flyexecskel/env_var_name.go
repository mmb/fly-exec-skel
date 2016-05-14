package flyexecskel

import (
	"strings"
)

func EnvVarName(s string) string {
	return strings.Replace(strings.ToUpper(s), "-", "_", -1)
}
