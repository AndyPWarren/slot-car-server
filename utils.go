package main

import (
	"fmt"
	"strings"
)

type KeyValue struct {
	key   string
	value string
}

func ParseKeyValueStr(str string, sep string) (error, KeyValue) {
	split := strings.Split(str, sep)
	if len(split) == 1 {
		return fmt.Errorf("string cant be parsed: %v does not exist in %v", sep, str), KeyValue{}
	} else {
		return nil, KeyValue{split[0], split[1]}
	}
}
