package pins

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	piblaster "github.com/AndyPWarren/go-pi-blaster"
)

type Pin struct {
	BcmPin int64
	Active bool
}

func NewPin(bcmPin int64) *Pin {
	return &Pin{bcmPin, false}
}

// AllPins is a map of configured pins where the key is the pin name (lane number) and the value is a pin object
var AllPins = make(map[string]*Pin)

var blaster = piblaster.Blaster{}

func square(val float64) float64 {
	return math.Pow(val, 2)
}

func cleanUp() {
	for _, pin := range AllPins {
		// fmt.Printf("cleaning up pin: %v \n", pin)
		blaster.Apply(pin.BcmPin, 0)
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
// It takes a map of pin names to pins,
// starts pi-blaster with these pins and starts the clean up watch tasks
func Setup(pins map[string]*Pin) {
	AllPins = pins
	bcmPins := make([]int64, len(pins))
	i := 0
	for _, val := range pins {
		bcmPins[i] = val.BcmPin
		i++
	}
	// fmt.Printf("starting pins: %v \n", bcmPins)
	blaster.Start(bcmPins)
	defer cleanUp()
	go watchForKill()
}

// Apply takes a pin name and a value and applies it to the pi-blaster, if the pin name has been configured. If it hasn't it returns an error
func Apply(pinName string, value float64) error {
	if pin, exists := AllPins[pinName]; exists == true {
		pin := pin.BcmPin
		// fmt.Printf("applying %v to pin: %v \n", value, pin)
		blaster.Apply(pin, square(value))
		return nil
	} else {
		return fmt.Errorf("pin not recognized: %v", pinName)
	}
}
