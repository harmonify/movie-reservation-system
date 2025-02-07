package model

import (
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
)

type UserOutbox struct {
	ID                 string    `json:"id"`
	Aggregatetype      string    `json:"aggregatetype"`
	Aggregateid        string    `json:"aggregateid"`
	Payload            []byte    `json:"payload"`
	Tracingspancontext []byte    `json:"tracingspancontext"`
	CreatedAt          time.Time `json:"created_at"`
}

func (m *UserOutbox) TableName() string {
	return "user_outbox"
}

func (m *UserOutbox) ToEntity() *entity.UserOutbox {
	return &entity.UserOutbox{
		ID:                 m.ID,
		AggregateType:      m.Aggregatetype,
		AggregateID:        m.Aggregateid,
		Payload:            m.Payload,
		Tracingspancontext: m.Tracingspancontext,
		CreatedAt:          m.CreatedAt,
	}
}

func NewUserOutbox(e entity.SaveUserOutbox) *UserOutbox {
	return &UserOutbox{
		ID:                 e.ID,
		Aggregatetype:      e.AggregateType,
		Aggregateid:        e.AggregateID,
		Payload:            e.Payload,
		Tracingspancontext: e.Tracingspancontext,
	}
}
