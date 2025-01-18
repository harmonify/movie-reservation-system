package test

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
)

type TestMessageListener struct {
	messages chan *kafka.ChanneledMessage
}

func (l *TestMessageListener) Messages() <-chan *kafka.ChanneledMessage {
	return l.messages
}

func (l *TestMessageListener) OnMessage(ctx context.Context, message *sarama.ConsumerMessage) {
	l.messages <- &kafka.ChanneledMessage{Context: ctx, Message: message}
}

func (l *TestMessageListener) Close() {
	close(l.messages)
}

func NewTestMessageListener() kafka.MessageListener {
	return &TestMessageListener{
		messages: make(chan *kafka.ChanneledMessage, 100),
	}
}
