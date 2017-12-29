package leds

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	piblaster "github.com/ddrager/go-pi-blaster"
)

type Message struct {
	color string
	value float64
}

const redLed = 17
const yellowLed = 27
const greenLed = 22

var leds = []int64{redLed, yellowLed, greenLed}

var blaster = piblaster.Blaster{}

func square(val float64) float64 {
	return math.Pow(val, 2)
}

func cleanUp() {
	for _, led := range leds {
		blaster.Apply(led, 0)
	}
}

func watchForKill() {
	killchan := make(chan os.Signal, 2)
	signal.Notify(killchan, os.Interrupt, syscall.SIGTERM)
	<-killchan
	log.Println("Kill sig!")
	cleanUp()
	os.Exit(0)
}

func Setup() {
	blaster.Start(leds)
	defer cleanUp()
	go watchForKill()
}

func Apply(inputColor string, value float64) error {
	var color int64
	switch inputColor {
	case "red":
		color = redLed
	case "yellow":
		color = yellowLed
	case "green":
		color = greenLed
	default:
		return fmt.Errorf("color not recognized: %v", inputColor)
	}
	blaster.Apply(color, value)
	return nil
}
