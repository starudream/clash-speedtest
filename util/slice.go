package util

import (
	"strings"
)

func Contains(vs []string, v string) bool {
	for _, t := range vs {
		if strings.Contains(v, t) {
			return true
		}
	}
	return false
}
