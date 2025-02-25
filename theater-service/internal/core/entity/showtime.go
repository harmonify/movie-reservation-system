package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Showtime struct {
	ShowtimeID string         `json:"showtime_id"`
	TraceID    string         `json:"trace_id"`
	TheaterID  string         `json:"theater_id"`
	RoomID     string         `json:"room_id"`
	MovieID    string         `json:"movie_id"`
	StartTime  time.Time      `json:"start_time"`
	EndTime    time.Time      `json:"end_time"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func (*Showtime) TableName() string {
	return "showtime"
}

func NewShowtime(save *SaveShowtime) *Showtime {
	return &Showtime{
		TraceID:   save.TraceID,
		RoomID:    save.RoomID,
		MovieID:   save.MovieID,
		StartTime: time.Time(save.StartTime),
		EndTime:   time.Time(save.EndTime),
	}
}

type FindOneShowtime struct {
	ShowtimeID sql.NullString
	TraceID    sql.NullString
}

type FindManyShowtimes struct {
	TheaterID    sql.NullString
	RoomID       sql.NullString
	MovieID      sql.NullString
	StartTimeGte sql.NullTime
	StartTimeLte sql.NullTime
	SortBy       ShowtimeSortBy
	Page         uint32
	PageSize     uint32
}

type ShowtimeSortBy string

const (
	ShowtimeSortByLatest ShowtimeSortBy = "latest"
	ShowtimeSortByOldest ShowtimeSortBy = "oldest"
)

type FindManyShowtimesResult struct {
	Showtimes []*Showtime
	Metadata  FindManyShowtimesMeta
}

type FindManyShowtimesMeta struct {
	TotalResults int64
}

type CountShowtimeTicket struct {
	ShowtimeID string
	Count      uint32
}

type SaveShowtime struct {
	TraceID   string
	RoomID    string
	MovieID   string
	StartTime time.Time
	EndTime   time.Time
}

type SaveShowtimeResult struct {
	ShowtimeID string
}

type UpdateShowtime struct {
	ShowtimeID sql.NullString
	TraceID    sql.NullString
	RoomID     sql.NullString
	MovieID    sql.NullString
	StartTime  sql.NullTime
	EndTime    sql.NullTime
}
