package chathub

type HubError struct {
	Code    string
	Message string
	Err     error
}

func (e *HubError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	} else {
		return e.Message
	}
}

var (
	ErrActiveClientWasNotFound = &HubError{
		Code:    "ACTIVE_CLIENT_NOT_FOUND",
		Message: "active client not found",
	}
)