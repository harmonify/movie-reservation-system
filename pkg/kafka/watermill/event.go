package watermill_pkg

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type (
	// Mainly used for testing purposes
	MessageListener interface {
		Messages() <-chan *ChanneledMessage
		OnMessage(ctx context.Context, message *message.Message)
		Close()
	}

	ChanneledMessage struct {
		Context context.Context
		Message *message.Message
	}
)
