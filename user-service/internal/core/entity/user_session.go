package entity

import (
	"database/sql"
	"time"
)

type (
	UserSession struct {
		ID           uint64
		UserUUID     string
		RefreshToken string // hashed
		IsRevoked    bool
		ExpiredAt    time.Time
		IpAddress    sql.NullString
		UserAgent    sql.NullString
		CreatedAt    time.Time
		UpdatedAt    time.Time
		DeletedAt    sql.NullTime
	}

	FindUserSession struct {
		ID           uint64
		UserUUID     sql.NullString
		RefreshToken sql.NullString
		IsRevoked    sql.NullBool
		ExpiredAt    sql.NullTime
		IpAddress    sql.NullString
		UserAgent    sql.NullString
		CreatedAt    sql.NullTime
		UpdatedAt    sql.NullTime
		DeletedAt    sql.NullTime
	}

	SaveUserSession struct {
		UserUUID     string
		RefreshToken string
		ExpiredAt    time.Time
		IpAddress    sql.NullString
		UserAgent    sql.NullString
	}

	UpdateUserSession struct {
		IsRevoked bool
		ExpiredAt time.Time
	}
)
