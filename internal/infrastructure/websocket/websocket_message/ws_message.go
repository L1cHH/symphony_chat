package websocketmessage

import (
	actions "symphony_chat/internal/domain/chat_actions"
)

type WsMessage struct {
	ChatAction actions.ChatActionType `json:"chat_action"`
	Payload    map[string]interface{} `json:"payload"`
}