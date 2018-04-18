package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"pi/slot-car-server/pins"
	"pi/slot-car-server/socket"
	"pi/slot-car-server/utils"
	"strconv"
)

func getPins(w http.ResponseWriter, r *http.Request) {
	pinNames := make([]string, len(pins.AllPins))
	i := 0
	for key := range pins.AllPins {
		pinNames[i] = key
		i++
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(pinNames)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, string(data))
}

func server() {
	http.HandleFunc("/lanes", getPins)
	http.HandleFunc("/", socket.Handler)     // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	args := os.Args[1:]
	allPins := make(map[string]*pins.Pin)
	for _, arg := range args {
		parsedArg, err := utils.ParseKeyValueStr(arg, "=")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			val, _ := strconv.ParseInt(parsedArg.Value, 0, 64)
			allPins[parsedArg.Key] = pins.NewPin(val)
		}
	}
	pins.Setup(allPins)
	server()
}
