package client

import (
	"encoding/json"
	"log"
	websocketmessage "symphony_chat/internal/infrastructure/websocket/websocket_message"
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
}

type Client struct {
	// conn is the websocket connection
	conn *websocket.Conn

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

func (c *Client) CloseConnection() {
	close(c.sendBuffer)
	close(c.receiveBuffer)
	c.conn.Close()
}

//Read messages from the connection and send them to the message receiver
func (c *Client) ProcessAndSendMessages() {
	for message := range c.receiveBuffer {
		c.msgReceiver.HandleMessage(message)
	}
}

//Handle messages from other clients
func (c *Client) HandleMessageFromServer(message []byte) {
	var msg websocketmessage.WsMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("error unmarshalling message: %v", err)
		return
	}
	//TODO: handle message
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.CloseConnection()
	}()
	
	for {
		select {
		case message := <-c.sendBuffer:
			go c.HandleMessageFromServer(message)
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

