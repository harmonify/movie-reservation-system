package test

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/kafka"
)

// Example of a custom listener for testing
type TestEventListener struct {
	events chan *kafka.ChanneledEvent
}

func (l *TestEventListener) Events() <-chan *kafka.ChanneledEvent {
	return l.events
}

func (l *TestEventListener) OnEvent(ctx context.Context, event *kafka.Event) {
	l.events <- &kafka.ChanneledEvent{Context: ctx, Event: event}
}

func (l *TestEventListener) Close() {
	close(l.events)
}

func NewTestEventListener() kafka.EventListener {
	return &TestEventListener{
		events: make(chan *kafka.ChanneledEvent, 100),
	}
}
