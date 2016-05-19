package flyexecskel

import (
	"fmt"
	"strings"
)

func Divider(label string) string {
	return fmt.Sprintf("# %s ", label) + strings.Repeat("-", 77-len(label))
}
