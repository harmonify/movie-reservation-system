package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusBooked    SeatStatus = "booked"
)

type Seat struct {
	SeatID     string         `json:"seat_id"`
	TraceID    string         `json:"trace_id"`
	RoomID     string         `json:"room_id"`
	SeatRow    string         `json:"seat_row"`
	SeatColumn string         `json:"seat_column"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func (*Seat) TableName() string {
	return "seat"
}

type FindManySeats struct {
	TraceID sql.NullString
	RoomID  sql.NullString
}

type FindOneSeat struct {
	SeatID sql.NullString
}

type CountRoomSeats struct {
	RoomID string
	Count  uint32
}

type FindShowtimeAvailableSeats struct {
	ShowtimeID string
}

type SaveSeat struct {
	TraceID string
	RoomID  string
	Row     string
	Column  string
}

type UpdateSeat struct {
	RoomID sql.NullString
	Row    sql.NullString
	Column sql.NullString
}
