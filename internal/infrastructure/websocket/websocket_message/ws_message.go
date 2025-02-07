package websocketmessage

import (
	actions "symphony_chat/internal/domain/chat_actions"
)

type WsMessageRequest struct {
	ChatAction actions.ChatActionType `json:"chat_action"`
	Payload    map[string]interface{} `json:"payload"`
}

type WsMessageResponse struct {
	ChatActionResult actions.ChatActionResult `json:"chat_action_result"`
	//If ChatActionResult == FAILED => Payload contsins "error" key
	Payload          map[string]interface{}   `json:"payload"`
}