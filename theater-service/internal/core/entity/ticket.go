package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	TicketID      string         `json:"ticket_id"`
	TraceID       string         `json:"trace_id"`
	TheaterID     string         `json:"theater_id"`
	RoomID        string         `json:"room_id"`
	SeatID        string         `json:"seat_id"`
	MovieID       string         `json:"movie_id"`
	ShowtimeID    string         `json:"showtime_id"`
	ReservationID string         `json:"reservation_id"`
	Price         float64        `json:"price"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

func (*Ticket) TableName() string {
	return "ticket"
}

type FindManyTickets struct {
	TheaterID  sql.NullString
	RoomID     sql.NullString
	ShowtimeID sql.NullString
}

type FindOneTicket struct {
	TicketID      sql.NullString
	TraceID       sql.NullString
	SeatID        sql.NullString
	ReservationID sql.NullString
}

type SaveTicket struct {
	TicketID      string
	TraceID       string
	TheaterID     string
	ShowtimeID    string
	SeatID        string
	ReservationID string
}

type UpdateTicket struct {
	TicketID      sql.NullString
	TraceID       sql.NullString
	TheaterID     sql.NullString
	ShowtimeID    sql.NullString
	SeatID        sql.NullString
	ReservationID sql.NullString
}
