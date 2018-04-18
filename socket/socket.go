package socket

import (
	"fmt"
	"log"
	"net/http"
	"pi/slot-car-server/pins"
	"pi/slot-car-server/utils"
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

	// clients pin given when they make a connection
	pin *pins.Pin
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
		fmt.Printf("error is %v\n", err)
		fmt.Printf("message is %v\n", message)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Printf("client closed the connection, release their pin, \n %v \n", err)
				c.pin.Active = false
			}
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v \n", err)
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
	var available = []string{}
	i := 0
	for key, pin := range pins.AllPins {
		if pin.Active != true {
			available = append(available, key)
		}
		i++
	}
	fmt.Printf("pins the are available: %v \n", available)
	if len(available) == 0 {
		fmt.Printf("no pins available \n")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256), pin: pins.AllPins[available[0]]}
	fmt.Printf("active pin is %v \n", available[0])
	client.send <- []byte(fmt.Sprintf("channel=%v", available[0]))
	activePin := pins.AllPins[available[0]]
	activePin.Active = true
	go client.read()
	go client.write()
}
