package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"symphony_chat/internal/infrastructure/websocket/chathub"
	"symphony_chat/internal/infrastructure/websocket/client"
	service "symphony_chat/internal/service/chat"
)

type WebsocketHandler struct {
	hub *chathub.Hub
	upgrader websocket.Upgrader
}

func NewWebsocketHandler(chatService *service.ChatService) *WebsocketHandler {
	return &WebsocketHandler {
		hub: chathub.NewHub(chatService),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
		},
	}
}


func (wh *WebsocketHandler) HandleWebSocket(c *gin.Context) {
	conn, err := wh.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "INTERNAL_SERVER_ERROR",
			"message": "internal server error, cannot upgrade http connection to websocket",
			"details": err.Error(),
		})
		return
	}

	userIdValue, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "INTERNAL_SERVER_ERROR",
			"message": "problem with getting user id from context",
			"details": "user id was not provided",
		})
		conn.Close()
		return
	}

	userID, ok := userIdValue.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "INTERNAL_SERVER_ERROR",
			"message": "problem with parsing user id",
			"details": "user id was not parsed to uuid",
		})
		conn.Close()
		return
	}

	client := client.NewClient(conn, userID, wh.hub)

	wh.hub.AddActiveClient(client)

	go client.ReadPump()
	go client.WritePump()
}

