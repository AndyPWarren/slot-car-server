package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pi/leds/leds"
	"pi/leds/socket"
	"strconv"
)

func server() {
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
