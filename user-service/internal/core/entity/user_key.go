package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type (
	UserKey struct {
		UserUUID   uuid.UUID
		PublicKey  string
		PrivateKey string `json:"-"` // encrypted
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  sql.NullTime
	}

	FindUserKey struct {
		UserUUID   sql.NullString
		PublicKey  sql.NullString
		// PrivateKey sql.NullString
		CreatedAt  sql.NullTime
		UpdatedAt  sql.NullTime
		DeletedAt  sql.NullTime
	}

	SaveUserKey struct {
		UserUUID   uuid.UUID
		PublicKey  string
		PrivateKey string `json:"-"` // encrypted
	}

	UpdateUserKey struct {
		UserUUID   uuid.UUID
		PublicKey  string
		PrivateKey string `json:"-"` // encrypted
	}
)
