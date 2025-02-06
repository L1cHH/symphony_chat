package client

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	MaxMessageSize = 1024

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	// conn is the websocket connection
	conn *websocket.Conn

	// sendBuffer is a channel for sending messages
	sendBuffer chan []byte

	// receiveBuffer is a channel for receiving messages
	receiveBuffer chan []byte

	userID uuid.UUID
}

func (c *Client) GetID() uuid.UUID {
	return c.userID
}

func NewClient(conn *websocket.Conn, userID uuid.UUID) *Client {
	if conn == nil {
		return nil
	}

	return &Client{
		conn: conn,
		sendBuffer: make(chan []byte),
		receiveBuffer: make(chan []byte),
		userID: userID,
	}
}

func (c *Client) CloseConnection() {
	close(c.sendBuffer)
	close(c.receiveBuffer)
	c.conn.Close()
}

//By client message handling means that the message was sent by the current user(current connection)
func (c *Client) HandleMessageFromClient(message []byte) {
	//TODO: handle client message
}

//By server message handling means that the message was sent by other clients(other users)
func (c *Client) HandleMessageFromServer(message []byte) {
	//TODO: handle server message
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
			c.HandleMessageFromServer(message)
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

		c.HandleMessageFromClient(message)
	}

}

