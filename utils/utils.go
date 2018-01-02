// package utils provides generic functions, light-weight functions which are useful across the project
package utils

import (
	"fmt"
	"strings"
)

// KeyValue defines a struct for key, value pairs
type KeyValue struct {
	Key   string
	Value string
}

// ParseKeyValueStr takes a string and a separator character
// and returns a KeyValue object of the string separator by the separator character
// if the separator character is not found its returns an error
func ParseKeyValueStr(str string, sep string) (KeyValue, error) {
	split := strings.Split(str, sep)
	if len(split) == 1 {
		return KeyValue{}, fmt.Errorf("string cant be parsed: %v does not exist in %v", sep, str)
	}
	return KeyValue{split[0], split[1]}, nil
}
