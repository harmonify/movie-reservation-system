package kafka

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

type (
	Event struct {
		Headers   []*sarama.RecordHeader `json:"headers"`
		Timestamp time.Time              `json:"timestamp,omitempty"`
		Key       string                 `json:"key,omitempty"`
		Value     []byte                 `json:"value,omitempty"`
		Topic     string                 `json:"topic,omitempty"`
	}

	EventListener interface {
		Events() <-chan *ChanneledEvent
		OnEvent(ctx context.Context, event *Event)
		Close()
	}

	ChanneledEvent struct {
		Context context.Context
		Event   *Event
	}

	// Mainly used for testing purposes
	MessageListener interface {
		Messages() <-chan *ChanneledMessage
		OnMessage(ctx context.Context, message *sarama.ConsumerMessage)
		Close()
	}

	ChanneledMessage struct {
		Context context.Context
		Message *sarama.ConsumerMessage
	}
)
