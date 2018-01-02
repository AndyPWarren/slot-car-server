package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pi/leds/leds"
	"pi/leds/socket"
	"pi/leds/utils"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	allLeds := make(map[string]int64)
	for _, arg := range args {
		parsedArg, err := utils.ParseKeyValueStr(arg, "=")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			val, _ := strconv.ParseInt(parsedArg.Value, 0, 64)
			allLeds[parsedArg.Key] = val
		}
	}
	leds.Setup(allLeds)
	server()
}
