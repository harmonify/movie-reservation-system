package model

import (
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"gorm.io/gorm"
)

type UserKey struct {
	UserUUID   string
	PublicKey  string
	PrivateKey string `json:"-"` // encrypted
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (m *UserKey) TableName() string {
	return "user_keys"
}

func (m *UserKey) ToEntity() *entity.UserKey {
	return &entity.UserKey{
		UserUUID:   m.UserUUID,
		PublicKey:  m.PublicKey,
		PrivateKey: m.PrivateKey,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  sql.NullTime(m.DeletedAt),
	}
}

func NewUserKey(e entity.SaveUserKey) *UserKey {
	return &UserKey{
		UserUUID:   e.UserUUID,
		PublicKey:  e.PublicKey,
		PrivateKey: e.PrivateKey,
	}
}
