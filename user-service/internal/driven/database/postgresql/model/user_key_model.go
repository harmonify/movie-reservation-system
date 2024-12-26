package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"gorm.io/gorm"
)

type UserKey struct {
	UserUUID   uuid.UUID
	PublicKey  string
	PrivateKey string `json:"-"` // encrypted
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserUUID"`
}

func (m *UserKey) TableName() string {
	return "user_keys"
}

func (m *UserKey) FromSaveEntity(e entity.SaveUserKey) *UserKey {
	return &UserKey{
		UserUUID:   e.UserUUID,
		PublicKey:  e.PublicKey,
		PrivateKey: e.PrivateKey,
	}
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
