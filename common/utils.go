package common

import "strings"

func ConcatStrings(params ...string) string {
	var result strings.Builder
	for _, s := range params {
		result.WriteString(s)
		result.WriteString(" ")
	}
	return result.String()
}
