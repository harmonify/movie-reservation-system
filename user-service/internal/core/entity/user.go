package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID                  uuid.UUID    `json:"uuid"`
	TraceID               string       `json:"trace_id"`
	Username              string       `json:"username"`
	Password              string       `json:"-"` // hashed
	Email                 string       `json:"email"`
	PhoneNumber           string       `json:"phone_number"`
	FirstName             string       `json:"first_name"`
	LastName              string       `json:"last_name"`
	IsEmailVerified       bool         `json:"is_email_verified"`
	IsPhoneNumberVerified bool         `json:"is_phone_number_verified"`
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`
	DeletedAt             sql.NullTime `json:"deleted_at"`
}

func (e *User) FullName() string {
	return e.FirstName + " " + e.LastName
}

type FindUser struct {
	UUID                  sql.NullString
	Username              sql.NullString
	Email                 sql.NullString
	PhoneNumber           sql.NullString
	FirstName             sql.NullString
	LastName              sql.NullString
	IsEmailVerified       sql.NullBool
	IsPhoneNumberVerified sql.NullBool
	CreatedAt             sql.NullTime
	UpdatedAt             sql.NullTime
	DeletedAt             sql.NullTime
}

type SaveUser struct {
	TraceID     string
	Username    string
	Password    string
	Email       string
	PhoneNumber string
	FirstName   string
	LastName    string
}

type UpdateUser struct {
	Email                 sql.NullString
	PhoneNumber           sql.NullString
	Username              sql.NullString
	FirstName             sql.NullString
	LastName              sql.NullString
	IsEmailVerified       sql.NullBool
	IsPhoneNumberVerified sql.NullBool
}
