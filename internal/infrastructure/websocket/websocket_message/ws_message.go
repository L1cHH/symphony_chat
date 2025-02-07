package websocketmessage

import (
	actions "symphony_chat/internal/domain/chat_actions"
)

//WebSocket message from client
type WsMessageRequest struct {
	ChatAction actions.ChatActionType `json:"chat_action"`
	Payload    map[string]interface{} `json:"payload"`
}

//WebSocket message to client (as response to WsMessageRequest)
type WsMessageResponse struct {
	ChatActionResult actions.ChatActionResult `json:"chat_action_result"`
	//If ChatActionResult == FAILED => Payload contsins "error" key
	Payload          map[string]interface{}   `json:"payload"`
}

//WebSocket client event (this event is used for sending events to engaged clients)
type WsClientEvent struct {
	EventType actions.EventType `json:"event_type"`
	Payload   map[string]interface{} `json:"payload"`
}