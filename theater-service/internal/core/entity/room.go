package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Room struct {
	RoomID    string         `json:"room_id"`
	TraceID   string         `json:"trace_id"`
	TheaterID string         `json:"theater_id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (*Room) TableName() string {
	return "room"
}

type FindOneRoom struct {
	RoomID  sql.NullString
	TraceID sql.NullString
}

type FindManyRooms struct {
	TheaterID sql.NullString
	Name      sql.NullString
}

type SaveRoom struct {
	TheaterID string
	TraceID   string
	Name      string
}

type UpdateRoom struct {
	TheaterID sql.NullString
	Name      sql.NullString
}
