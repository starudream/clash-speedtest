package util

import (
	"strings"
)

func StringContains(s string, ss ...string) bool {
	for _, v := range ss {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}
