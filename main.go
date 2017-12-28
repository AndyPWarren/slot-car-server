package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ddrager/go-pi-blaster"
)

const redLed = 17
const yellowLed = 27
const greenLed = 22

var leds = []int64{redLed, yellowLed, greenLed}

var blaster = piblaster.Blaster{}

func cleanUp() {
	for _, led := range leds {
		blaster.Apply(led, 0)
	}
}

func square(val float64) float64 {
	return math.Pow(val, 2)
}

func watchForKill() {
	killchan := make(chan os.Signal, 2)
	signal.Notify(killchan, os.Interrupt, syscall.SIGTERM)
	<-killchan
	log.Println("Kill sig!")
	cleanUp()
	os.Exit(0)
}

func main() {
	blaster.Start(leds)
	defer cleanUp()
	go watchForKill()
	fmt.Printf("Running\n")

	for i := 0; i < 11; i++ {
		step := float64(i) / 10
		blaster.Apply(redLed, square(step))
		blaster.Apply(yellowLed, step)
		time.Sleep(time.Second)
	}
}
