package sqlrepo

import "strings"

//List returns a comma separated list of value repeated count times.
//A non-positive count results in the empty string.
func List(value string, count int) string {
	if count <= 0 {
		return ""
	}

	result := strings.Repeat(value+",", count-1)
	result += value

	return result
}
