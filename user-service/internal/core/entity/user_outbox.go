package entity

import "time"

const (
	AggregateTypeRegistered = "registered"
)

type UserOutbox struct {
	ID                 string    `json:"id"`
	AggregateType      string    `json:"aggregatetype"`      // Aggregate event type, e.g. registered
	AggregateID        string    `json:"aggregateid"`        // a.k.a. user ID
	Payload            []byte    `json:"payload"`            // Protobuf binary data
	Tracingspancontext []byte    `json:"tracingspancontext"` // JSON binary data
	CreatedAt          time.Time `json:"created_at"`
}

type SaveUserOutbox struct {
	ID                 string `json:"id"`
	AggregateType      string `json:"aggregatetype"`
	AggregateID        string `json:"aggregateid"`
	Payload            []byte `json:"payload"`
	Tracingspancontext []byte `json:"tracingspancontext"`
}
