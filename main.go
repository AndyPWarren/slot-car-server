package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pi/leds/leds"
	"pi/leds/socket"
	"strconv"
)

func getLeds(w http.ResponseWriter, r *http.Request) {
	colors := make([]string, len(leds.AllLeds))
	i := 0
	for led := range leds.AllLeds {
		colors[i] = led
		i++
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(colors)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, string(data))
}

func server() {
	http.HandleFunc("/leds", getLeds)
	http.HandleFunc("/", socket.Handler)     // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	args := os.Args[1:]
	allLeds := make([]leds.Led, len(args))
	for i, arg := range args {
		err, parsedArg := ParseKeyValueStr(arg, "=")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			val, _ := strconv.ParseInt(parsedArg.value, 0, 64)
			allLeds[i] = leds.Led{parsedArg.key, val}
		}
	}
	leds.Setup(allLeds)
	server()
}
