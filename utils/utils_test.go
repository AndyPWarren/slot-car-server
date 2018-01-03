package utils

import (
	"testing"
)

type TestCase struct {
	name     string
	str      string
	sep      string
	expected KeyValue
}

func TestParseKeyValueStr(t *testing.T) {
	cases := []TestCase{
		TestCase{
			"equals is the separator character and exists",
			"test=case",
			"=",
			KeyValue{"test", "case"}},
		TestCase{"colon is the separator character and exists", "test:case", ":", KeyValue{"test", "case"}},
		TestCase{"equals is the separator character and does not exist in string", "testcase", "=", KeyValue{}},
	}
	for i, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			result, err := ParseKeyValueStr(test.str, test.sep)
			if result != test.expected {
				t.Error(test, result)
			}
			if i == 2 {
				if err == nil {
					t.Error(test, result)
				}
			}
		})
	}

}
