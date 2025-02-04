package client

import (
	"encoding/json"
	websocketmessage "symphony_chat/internal/infrastructure/websocket/websocket_message"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

func (c *Client) CloseBufferChannels() {
	close(c.sendBuffer)
	close(c.receiveBuffer)
}

func (c *Client) HandleMessage(message []byte) error {
	var msg websocketmessage.WsMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return err
	}

	return nil
}