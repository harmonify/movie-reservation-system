package entity

import (
	"database/sql"
	"time"
)

type User struct {
	UUID                  string       `json:"uuid" faker:"uuid_hyphenated"`
	TraceID               string       `json:"trace_id" faker:"uuid_hyphenated"`
	Username              string       `json:"username" faker:"username"`
	Password              string       `json:"-" faker:"-"` // hashed
	Email                 string       `json:"email" faker:"email"`
	PhoneNumber           string       `json:"phone_number" faker:"e_164_phone_number"`
	FirstName             string       `json:"first_name" faker:"first_name"`
	LastName              string       `json:"last_name" faker:"last_name"`
	IsEmailVerified       bool         `json:"is_email_verified"`
	IsPhoneNumberVerified bool         `json:"is_phone_number_verified"`
	CreatedAt             time.Time    `json:"created_at" faker:"-"`
	UpdatedAt             time.Time    `json:"updated_at" faker:"-"`
	DeletedAt             sql.NullTime `json:"deleted_at" faker:"-"`
}

func (e *User) FullName() string {
	return e.FirstName + " " + e.LastName
}

type GetUser struct {
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
	Password    string `json:"-"`
	Email       string
	PhoneNumber string
	FirstName   string
	LastName    string
}

type UpdateUser struct {
	Username              sql.NullString
	Password              sql.NullString `json:"-"`
	Email                 sql.NullString
	PhoneNumber           sql.NullString
	FirstName             sql.NullString
	LastName              sql.NullString
	IsEmailVerified       sql.NullBool
	IsPhoneNumberVerified sql.NullBool
}

type UserEmail struct {
	Email string
}

type UserPhoneNumber struct {
	PhoneNumber string
}
