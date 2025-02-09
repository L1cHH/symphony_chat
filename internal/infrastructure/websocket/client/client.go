package client

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	//Max size of the connection's incoming messages
	MaxMessageSize = 1024

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	//Size of the buffers for sending and receiving messages
	sendBufferSize = 256
	receiveBufferSize = 256
)

type MessageReceiver interface {
	HandleMessage(message []byte)
	RemoveActiveClient(client *Client)
}

type Client struct {
	// conn is the websocket connection
	conn *websocket.Conn

	// closeOnce ensures that CloseConnection is called at most once
	closeOnce sync.Once

	// sendBuffer is a channel for receiving messages from Hub
	sendBuffer chan []byte

	// receiveBuffer is a channel for receiving messages from current connection
	receiveBuffer chan []byte

	// userID defines a user of the current connection
	userID uuid.UUID

	// msgReceiver is a receiver for messages from current Client
	msgReceiver MessageReceiver
}

func (c *Client) GetID() uuid.UUID {
	return c.userID
}

func NewClient(conn *websocket.Conn, userID uuid.UUID, msgReceiver MessageReceiver) *Client {
	if conn == nil {
		return nil
	}

	return &Client{
		conn: conn,
		sendBuffer: make(chan []byte, sendBufferSize),
		receiveBuffer: make(chan []byte, receiveBufferSize),
		userID: userID,
		msgReceiver: msgReceiver,
	}
}

func (c *Client) IsStillConnected() bool {
	if c.conn == nil {
		return false
	} else {
		return true
	}
}

func (c *Client) CloseConnection() {
	c.closeOnce.Do(func() {
		close(c.sendBuffer)
		close(c.receiveBuffer)
		c.conn.Close()
		c.msgReceiver.RemoveActiveClient(c)
	})
} 

func (c *Client) GetMessageFromServer(message []byte) {
	c.sendBuffer <- message
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.CloseConnection()
	}()
	
	for {
		select {
		case message, ok := <-c.sendBuffer:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(pongWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}

}

func (c *Client) ReadPump() {
	defer func() {
		c.CloseConnection()
	}()

	c.conn.SetReadLimit(MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go c.ProcessAndSendMessages()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {log.Printf("error: %v", err)}

			break
		}
		c.receiveBuffer <- message
	}
}

//Read messages from the connection and send them to the message receiver
func (c *Client) ProcessAndSendMessages() {
	for message := range c.receiveBuffer {
		c.msgReceiver.HandleMessage(message)
	}
}

