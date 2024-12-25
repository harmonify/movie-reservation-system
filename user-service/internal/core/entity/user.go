package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID                  uuid.UUID
	Username              string
	Password              string `json:"-"` // hashed
	Email                 string
	PhoneNumber           string
	FirstName             string
	LastName              string
	IsEmailVerified       bool
	IsPhoneNumberVerified bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             sql.NullTime
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
	FirstName             sql.NullString
	LastName              sql.NullString
	IsEmailVerified       sql.NullBool
	IsPhoneNumberVerified sql.NullBool
}
