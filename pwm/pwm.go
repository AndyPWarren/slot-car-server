package pwm

import (
	"fmt"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

func HighDuration(freq float32, dc float32) time.Duration {
	val := freq * dc
	return time.Duration(val) * time.Millisecond
}

func LowDuration(freq float32, dc float32) time.Duration {
	val := freq * (1 - dc)
	return time.Duration(val) * time.Millisecond
}

type Pwm struct {
	Pin     rpio.Pin
	Freq    int
	Dc      chan float32
	Running chan bool
}

func run(pwm *Pwm) {
	fmt.Printf("running: %v", pwm.Pin)
	freqMs := float32(1000 / pwm.Freq)
	high := HighDuration(freqMs, 0)
	low := LowDuration(freqMs, 0)
	for {
		select {
		case r := <-pwm.Running:
			if r == false {
				fmt.Printf("stopping pin no: %v \n", pwm.Pin)
				pwm.Pin.Low()
				return
			}
		case dc := <-pwm.Dc:
			fmt.Printf("duty cycle: %v", dc)
			high = HighDuration(freqMs, dc)
			low = LowDuration(freqMs, dc)
		default:
			pwm.Pin.High()
			time.Sleep(high)
			pwm.Pin.Low()
			time.Sleep(low)
		}
	}
}

func Create(pin rpio.Pin, freq int, dc float32) *Pwm {
	runningCh := make(chan bool)
	dcChan := make(chan float32, 1)
	dcChan <- dc
	p := &Pwm{pin, freq, dcChan, runningCh}
	return p
}

func (p *Pwm) Start() {
	fmt.Printf("starting: %v", p.Pin)
	go run(p)

	p.Running <- true
}

func (p *Pwm) Stop() {
	p.Running <- false
}

func (p *Pwm) ChangeDutyCycle(dc float32) {
	fmt.Printf("ChangeDutyCycle: %v", dc)
	p.Dc <- dc
}
