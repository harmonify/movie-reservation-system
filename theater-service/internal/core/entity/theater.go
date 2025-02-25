package entity

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Theater struct {
	TheaterID   string         `json:"theater_id"`
	TraceID     string         `json:"trace_id"`
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	PhoneNumber string         `json:"phone_number"`
	Email       string         `json:"email"`
	Website     string         `json:"website"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (*Theater) TableName() string {
	return "theater"
}

type SaveTheater struct {
	TraceID     string
	Name        string
	Address     string
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Website     string `json:"website"`
}

type FindOneTheater struct {
	TheaterID   sql.NullString
	TraceID     sql.NullString
	PhoneNumber sql.NullString `json:"phone_number"`
	Email       sql.NullString `json:"email"`
	Website     sql.NullString `json:"website"`
}

type FindManyTheaters struct {
	Name    sql.NullString
	Address sql.NullString
}

type UpdateTheater struct {
	Name    sql.NullString
	Address sql.NullString
}
