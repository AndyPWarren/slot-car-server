package leds

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	piblaster "github.com/AndyPWarren/go-pi-blaster"
)

// AllLeds is a map of configured leds where the key is the color of the led and int64 is the pin value
// for that color
var AllLeds = make(map[string]int64)

var blaster = piblaster.Blaster{}

func square(val float64) float64 {
	return math.Pow(val, 2)
}

func cleanUp() {
	for _, pin := range AllLeds {
		blaster.Apply(pin, 0)
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

// Setup performs initial pi-blaster setup.
// It takes a map of led color to pin number,
// starts pi-blaster with these pins and starts the clean up watch tasks
func Setup(leds map[string]int64) {
	AllLeds = leds
	pins := make([]int64, len(leds))
	i := 0
	for _, val := range leds {
		pins[i] = val
		i++
	}
	blaster.Start(pins)
	defer cleanUp()
	go watchForKill()
}

// Apply takes a color and a brightness value and applies it to the pi-blaster, if the input color has been configured. If it hasn't it returns an error
func Apply(inputColor string, value float64) error {
	var pin int64
	if AllLeds[inputColor] != 0 {
		pin = AllLeds[inputColor]
	} else {
		return fmt.Errorf("color not recognized: %v", inputColor)
	}
	blaster.Apply(pin, square(value))
	return nil
}
