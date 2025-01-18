package kafka

import "errors"

var (
	// Producer errors below
	ErrNilMessage = errors.New("cannot send nil message")
	// Consumer errors below
	ErrMessageChannelClosed = errors.New("message channel was closed")
	ErrDecodeFailed         = errors.New("failed to decode event value")
	ErrInvalidValueType     = errors.New("invalid event value type")
)
