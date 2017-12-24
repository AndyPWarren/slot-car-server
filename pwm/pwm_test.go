package pwm

import (
	"testing"
	"time"
)

type DurationTestCase struct {
	freq     float32
	dc       float32
	expected time.Duration
}

func NewDurationTestCase(freq int, dc float32, expect int) DurationTestCase {
	return DurationTestCase{float32(1000 / freq), dc, time.Duration(expect) * time.Millisecond}
}

func TestHighDuration(t *testing.T) {
	case1 := NewDurationTestCase(100, 0.5, 5)
	case2 := NewDurationTestCase(1, 0.9, 900)
	case3 := NewDurationTestCase(100, 0.1, 1)
	tests := [3]DurationTestCase{case1, case2, case3}
	for _, test := range tests {
		if result := HighDuration(test.freq, test.dc); result != test.expected {
			t.Error(test, result)
		}
	}
}

func TestLowDuration(t *testing.T) {
	case1 := NewDurationTestCase(100, 0.5, 5)
	case2 := NewDurationTestCase(1, 0.9, 100)
	case3 := NewDurationTestCase(100, 0.1, 9)
	tests := [3]DurationTestCase{case1, case2, case3}
	for _, test := range tests {
		if result := LowDuration(test.freq, test.dc); result != test.expected {
			t.Error(test, result)
		}
	}
}
