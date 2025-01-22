package kafka

import (
	"errors"
)

var (
	// Producer errors below
	ErrNilMessage = errors.New("cannot send nil message")
	// Consumer errors below
	ErrMessageChannelClosed = errors.New("message channel was closed")
	ErrMalformedMessage     = errors.New("malformed message")
)
