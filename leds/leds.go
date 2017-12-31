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

type Led struct {
	Color string
	Pin   int64
}

var AllLeds = make(map[string]int64)

var blaster = piblaster.Blaster{}

func square(val float64) float64 {
	return math.Pow(val, 2)
}

func cleanUp() {
	for _, led := range AllLeds {
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

func Setup(leds []Led) {
	pins := make([]int64, len(leds))
	for i, led := range leds {
		AllLeds[led.Color] = led.Pin
		pins[i] = led.Pin
	}
	blaster.Start(pins)
	defer cleanUp()
	go watchForKill()
}

func Apply(inputColor string, value float64) error {
	var pin int64
	if AllLeds[inputColor] != 0 {
		pin = AllLeds[inputColor]
	} else {
		return fmt.Errorf("color not recognized: %v", inputColor)
	}
	blaster.Apply(pin, value)
	return nil
}
