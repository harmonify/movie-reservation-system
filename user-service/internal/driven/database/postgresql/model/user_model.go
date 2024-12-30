package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"gorm.io/gorm"
)

type User struct {
	UUID                  uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // For pg14+
	Username              string    `gorm:"uniqueIndex:idx_user_username;unique"`
	Password              string    `json:"-"`
	Email                 string    `gorm:"uniqueIndex:idx_user_email;unique"`
	PhoneNumber           string    `gorm:"uniqueIndex:idx_user_phone_number;unique"`
	FirstName             string
	LastName              string
	IsEmailVerified       bool `gorm:"default:false"`
	IsPhoneNumberVerified bool `gorm:"default:false"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`

	UserSessions []UserSession
}

func (m *User) TableName() string {
	return "users"
}

func (m *User) FromEntity(e entity.User) *User {
	return &User{
		UUID:                  e.UUID,
		Username:              e.Username,
		Email:                 e.Email,
		PhoneNumber:           e.PhoneNumber,
		FirstName:             e.FirstName,
		LastName:              e.LastName,
		IsEmailVerified:       e.IsEmailVerified,
		IsPhoneNumberVerified: e.IsPhoneNumberVerified,
		CreatedAt:             e.CreatedAt,
		UpdatedAt:             e.UpdatedAt,
		DeletedAt:             gorm.DeletedAt(e.DeletedAt),
	}
}

func (m *User) ToEntity() *entity.User {
	return &entity.User{
		UUID:                  m.UUID,
		Username:              m.Username,
		Password:              m.Password,
		Email:                 m.Email,
		PhoneNumber:           m.PhoneNumber,
		FirstName:             m.FirstName,
		LastName:              m.LastName,
		IsEmailVerified:       m.IsEmailVerified,
		IsPhoneNumberVerified: m.IsPhoneNumberVerified,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		DeletedAt:             sql.NullTime(m.DeletedAt),
	}
}

func (m *User) FromSaveEntity(e entity.SaveUser) *User {
	return &User{
		Username:              e.Username,
		Password:              e.Password,
		Email:                 e.Email,
		PhoneNumber:           e.PhoneNumber,
		FirstName:             e.FirstName,
		LastName:              e.LastName,
		IsEmailVerified:       false,
		IsPhoneNumberVerified: false,
	}
}

func (m *User) ToSaveEntity() *entity.SaveUser {
	return &entity.SaveUser{
		Username:    m.Username,
		Password:    m.Password,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
	}
}
