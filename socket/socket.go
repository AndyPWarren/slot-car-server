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

func WriteMessage(m string) {
	byteMessage := []byte(m)
	if err := conn.WriteMessage(1, byteMessage); err != nil {
		log.Println(err)
		return
	}
}

func parseMessage(m string) {
	err, parsedM := utils.ParseKeyValueStr(m, "=")
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
		} else {
			parseMessage(string(p))
		}

	}
}
