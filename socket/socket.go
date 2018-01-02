package socket

import (
	"fmt"
	"log"
	"net/http"
	"pi/leds/leds"
	"pi/leds/utils"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var conn *websocket.Conn

func parseMessage(m string) {
	parsedM, err := utils.ParseKeyValueStr(m, "=")
	if err != nil {
		WriteMessage(err.Error())
		return
	}
	color := parsedM.Key
	value, err := strconv.ParseFloat(parsedM.Value, 64)
	if err != nil {
		WriteMessage(fmt.Sprintf("error converting %v to float: %v", value, err))
		return
	}
	applyErr := leds.Apply(color, value)
	if applyErr != nil {
		WriteMessage(applyErr.Error())
	}
}

// WriteMessage takes a string message (m) and sends it onto the socket connection
func WriteMessage(m string) {
	byteMessage := []byte(m)
	if err := conn.WriteMessage(1, byteMessage); err != nil {
		log.Println(err)
		return
	}
}

// Handler takes a ResponseWriter and request and upgrades it to a socket connection
// It reads the incoming message and parses it
func Handler(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn = connection
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		parseMessage(string(p))
	}
}
