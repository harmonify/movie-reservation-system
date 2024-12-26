package messaging

import "context"

type Messager interface {
	Send(ctx context.Context, message Message) (id string, err error)
}

type Message struct {
	To   string // phone number
	Body string
}
