package netw

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type IBaseClient interface {
	ReadPump()
	WritePump()
}

// Client is a middleman between the websocket Connection and the hub.
type BaseClient struct {
	IBaseClient
	Conn       *websocket.Conn
	Send       chan *Envelope
	UserId     string
	SentBy     interface{}
	Notify     chan *Notify
	Unregister chan interface{}
}

func NewBaseClient(extended interface{}) *BaseClient {
	return &BaseClient{
		Send:   make(chan *Envelope, 10),
		Notify: make(chan *Notify, 1),
		SentBy: extended,
	}
}

// // readPump pumps messages from the websocket Connection to the hub.

// // The application runs readPump in a per-Connection goroutine. The application
// // ensures that there is at most one reader on a Connection by executing all
// // reads from this goroutine.
func (c *BaseClient) ReadPump() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
		fmt.Println("ReadPump defered")
		if c.Unregister != nil {
			c.Unregister <- c.SentBy
		}
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		log.Println(string(message[:]))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		var msg json.RawMessage
		env := Envelope{
			Message: &msg,
		}
		if err := json.Unmarshal(message, &env); err != nil {
			log.Fatal(err)
		}
		if c.Notify != nil {
			switch env.MessageCode {
			case EMessage:
				var message Message
				if err := json.Unmarshal(msg, &message); err != nil {
					log.Fatal(err)
				}
				env.Message = message
			case EEvent:
				var event Event
				if err := json.Unmarshal(msg, &event); err != nil {
					log.Fatal(err)
				}
				env.Message = event

			case ERegister:
				var register Register
				if err := json.Unmarshal(msg, &register); err != nil {
					log.Fatal(err)
				}
				env.Message = register
			default:
				continue
			}

			c.Notify <- &Notify{
				Message: &env,
				SentBy:  c.SentBy,
			}
		}

	}
}

// // writePump pumps messages from the hub to the websocket Connection.
// //
// // A goroutine running writePump is started for each Connection. The
// // application ensures that there is at most one writer to a Connection by
// // executing all writes from this goroutine.
func (c *BaseClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//TextMessage denotes a text data message.
			// The text message payload is interpreted as UTF-8 encoded text data.
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			msgToClient, _ := json.Marshal(message)
			// msg := &Envelope{
			// 	Client:      message.Client,
			// 	MessageCode: message.MessageCode,
			// 	Message:     message.Message,
			// }
			// msgToClient, _ = json.Marshal(msg)
			log.Println(string(msgToClient[:]))
			w.Write(msgToClient)
			if err := w.Close(); err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func (b *BaseClient) ServeWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error(err)
	}
	b.Conn = Conn
	go b.WritePump()
	go b.ReadPump()
}
