package main

import (
	"log"
	"net/http"
	"pi/leds/leds"
	"pi/leds/socket"
)

func server() {
	http.HandleFunc("/", socket.Handler)     // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	leds.Setup()
	server()
}
