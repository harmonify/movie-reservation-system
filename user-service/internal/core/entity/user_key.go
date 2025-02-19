package entity

import (
	"database/sql"
	"time"
)

type (
	UserKey struct {
		UserUUID   string       `json:"user_uuid"`
		PublicKey  string       `json:"public_key"`
		PrivateKey string       `json:"-"` // encrypted
		CreatedAt  time.Time    `json:"created_at"`
		UpdatedAt  time.Time    `json:"updated_at"`
		DeletedAt  sql.NullTime `json:"deleted_at"`
	}

	GetUserKey struct {
		UserUUID  sql.NullString
		PublicKey sql.NullString
		// PrivateKey sql.NullString
		CreatedAt sql.NullTime
		UpdatedAt sql.NullTime
		DeletedAt sql.NullTime
	}

	SaveUserKey struct {
		UserUUID   string
		PublicKey  string
		PrivateKey string `json:"-"` // encrypted
	}

	UpdateUserKey struct {
		UserUUID   string
		PublicKey  string
		PrivateKey string `json:"-"` // encrypted
	}
)
