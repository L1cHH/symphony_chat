package chat 

type ChatError struct {
	Code    string
	Message string
	Err     error
}

func (e *ChatError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	} else {
		return e.Message
	}
}

var (
	ErrWrongChatName = &ChatError {
		Code: "WRONG_CHAT_NAME",
		Message: "wrong chat name",
	}
)