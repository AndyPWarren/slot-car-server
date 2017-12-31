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
