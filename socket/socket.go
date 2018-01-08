package socket

import (
	"fmt"
	"log"
	"net/http"
	"pi/leds/pins"
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

type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func parseMessage(m string) error {
	parsedM, err := utils.ParseKeyValueStr(m, "=")
	if err != nil {
		return err
	}
	color := parsedM.Key
	value, parsingErr := strconv.ParseFloat(parsedM.Value, 64)
	if parsingErr != nil {
		return fmt.Errorf("error converting %v to float: %v", value, parsingErr)
	}
	applyErr := pins.Apply(color, value)
	if applyErr != nil {
		return applyErr
	}
	return nil
}

// WriteMessage takes a string message (m) and sends it onto the socket connection
func WriteMessage(m string, conn *websocket.Conn) {
	byteMessage := []byte(m)
	if err := conn.WriteMessage(1, byteMessage); err != nil {
		log.Println(err)
		return
	}
}
func (c *Client) write() {
	for {
		select {
		case message := <-c.send:
			if err := c.conn.WriteMessage(1, message); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Client) read() {
	defer func() {
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		// formatted := fmt.Sprintf("message received: %v \n", string(message))
		parseMessageError := parseMessage(string(message))
		fmt.Println(parseMessageError)
		if parseMessageError != nil {
			c.send <- []byte(parseMessageError.Error())
		}
	}
}

// Handler takes a ResponseWriter and request and upgrades it to a socket connection
// It reads the incoming message and parses it
func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	go client.read()
	go client.write()
}
