package core

var MESSAGE_OK int = 1
var MESSAGE_ERROR int = 0

type Message struct {
	Code    int
	Payload interface{}
}

func NewMessage(code int, payload interface{}) *Message {
	message := &Message{
		Code:    code,
		Payload: payload,
	}
	return message
}
