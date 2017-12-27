package pwm

import (
	"testing"
	"time"
)

type DurationTestCase struct {
	name     string
	freq     float32
	dc       float32
	expected time.Duration
}

func NewDurationTestCase(name string, freq int, dc float32, expect int) DurationTestCase {
	return DurationTestCase{name, float32(1000 / freq), dc, time.Duration(expect) * time.Millisecond}
}

func TestHighDuration(t *testing.T) {
	case1 := NewDurationTestCase("equal duty cycle", 100, 0.5, 5)
	case2 := NewDurationTestCase("low frequency", 1, 0.9, 900)
	case3 := NewDurationTestCase("short duty cycle", 100, 0.1, 1)
	tests := [3]DurationTestCase{case1, case2, case3}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := highDuration(test.freq, test.dc); result != test.expected {
				t.Error(test, result)
			}
		})
	}
}

func BenchmarkHighDuration(b *testing.B) {
	case1 := NewDurationTestCase("large duty cycle", 100, 0.9, 900)
	tests := [1]DurationTestCase{case1}
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				highDuration(test.freq, test.dc)
			}
		})
	}
}

func TestLowDuration(t *testing.T) {
	case1 := NewDurationTestCase("equal duty cycle", 100, 0.5, 5)
	case2 := NewDurationTestCase("low frequency", 1, 0.9, 100)
	case3 := NewDurationTestCase("short duty cycle", 100, 0.1, 9)
	tests := [3]DurationTestCase{case1, case2, case3}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := lowDuration(test.freq, test.dc); result != test.expected {
				t.Error(test, result)
			}
		})
	}
}
