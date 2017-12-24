package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"pi/leds/pwm"

	"github.com/stianeikeland/go-rpio"
)

type void func()

const led1 = rpio.Pin(17)
const led2 = rpio.Pin(27)
const led3 = rpio.Pin(22)

var leds = [3]rpio.Pin{led1, led2, led3}

func cleanUp() {
	for _, led := range leds {
		led.Low()
	}
	rpio.Close()
}

func ioSetUp() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("setting up...")
	for _, led := range leds {
		led.Output()
		led.Low()
	}
}

func kill(fn void) {
	killchan := make(chan os.Signal, 2)
	signal.Notify(killchan, os.Interrupt, syscall.SIGTERM)
	// wait for kill signal
	<-killchan
	log.Println("Kill sig!")
	//do clean up
	fn()
	//now exit
	os.Exit(0)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("please provide the duty cycle (number between 0 - 1) optional pass the frequency i(hz) e.g. 0.4 100")
		os.Exit(1)
	}
	dutyCycle, err := strconv.ParseFloat(os.Args[1], 32)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var freq int64 = 100
	if len(os.Args) == 2 {
		fmt.Printf("setting frequency to default value of %v", freq)
	} else {
		argFreg, err := strconv.ParseInt(os.Args[2], 10, 64)
		freq = argFreg
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Println("starting...")
	ioSetUp()
	defer cleanUp()
	pwm1 := pwm.Create(led1, int(freq), float32(dutyCycle))
	pwm2 := pwm.Create(led2, int(freq), float32(dutyCycle))
	pwm1.Start()
	pwm2.Start()
	var steps = 6
	for i := 0; i < steps; i++ {
		step := float32(i) / float32(steps-1)
		fmt.Println(step)
		pwm1.ChangeDutyCycle(step / 2)
		pwm2.ChangeDutyCycle(0.5 - (step / 2))
		time.Sleep(time.Second * 1)
	}
}
