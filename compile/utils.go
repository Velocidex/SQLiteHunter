package compile

import (
	"strings"
)

func indent(in string, indent int) string {
	lines := strings.Split(in, "\n")
	result := []string{}
	for _, l := range lines {
		result = append(result, strings.Repeat(" ", indent)+l)
	}
	return strings.Join(result, "\n")
}
