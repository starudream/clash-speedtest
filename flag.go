package main

import (
	"strings"
)

type stringsValue []string

func (s *stringsValue) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func (s *stringsValue) String() string { return strings.Join(*s, ",") }
