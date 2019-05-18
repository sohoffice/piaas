package stringarrays

import (
	"bytes"
	"fmt"
)

// Find the position of target in a string array
// Return -1 if not found
func IndexOf(ar []string, target string) int {
	for i, s := range ar {
		if s == target {
			return i
		}
	}
	return -1
}

func Compare(ar []string, that []string) bool {
	if len(ar) != len(that) {
		return false
	}
	for i := 0; i < len(ar); i++ {
		if (ar)[i] != (that)[i] {
			return false
		}
	}
	return true
}

// Reverse the string array
func Reverse(ar []string) []string {
	var target []string
	for i := range ar {
		target = append(target, (ar)[len(ar)-i-1])
	}
	return target
}

func ToString(ar []string) string {
	var buf bytes.Buffer
	for i, s := range ar {
		buf.WriteString(fmt.Sprintf("%2d. ", i))
		buf.WriteString(s)
		buf.WriteRune('\n')
	}
	return buf.String()
}
