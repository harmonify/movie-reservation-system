package entity

import (
	"database/sql"
	"time"
)

type (
	UserSession struct {
		ID           uint64         `json:"id"`
		UserUUID     string         `json:"user_uuid"`
		RefreshToken string         `json:"refresh_token"` // hashed
		IsRevoked    bool           `json:"is_revoked"`
		ExpiredAt    time.Time      `json:"expired_at"`
		IpAddress    sql.NullString `json:"ip_address"`
		UserAgent    sql.NullString `json:"user_agent"`
		CreatedAt    time.Time      `json:"created_at"`
		UpdatedAt    time.Time      `json:"updated_at"`
		DeletedAt    sql.NullTime   `json:"deleted_at"`
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
