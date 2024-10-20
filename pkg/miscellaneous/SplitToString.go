package miscellaneous

import (
	"strconv"
	"strings"
)

// SplitToString transform a slice of int to a string with a separator
// ex []int{1,2,3} -> "1,2,3"
func SplitToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
