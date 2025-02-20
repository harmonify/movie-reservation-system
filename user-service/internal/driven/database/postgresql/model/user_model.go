package model

import (
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"gorm.io/gorm"
)

const UserTableName = "users"

type User struct {
	UUID                  string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" faker:"-"` // For pg14+
	TraceID               string         `gorm:"uniqueIndex:idx_user_trace_id;unique" faker:"uuid_hyphenated"`
	Username              string         `gorm:"uniqueIndex:idx_user_username;unique" faker:"username"`
	Password              string         `json:"-" faker:"-"`
	Email                 string         `gorm:"uniqueIndex:idx_user_email;unique" faker:"email"`
	PhoneNumber           string         `gorm:"uniqueIndex:idx_user_phone_number;unique" faker:"phone_number"`
	FirstName             string         `faker:"first_name"`
	LastName              string         `faker:"last_name"`
	IsEmailVerified       bool           `gorm:"default:false" faker:"bool"`
	IsPhoneNumberVerified bool           `gorm:"default:false" faker:"bool"`
	CreatedAt             time.Time      `gorm:"autoCreateTime" faker:"-"`
	UpdatedAt             time.Time      `gorm:"autoUpdateTime" faker:"-"`
	DeletedAt             gorm.DeletedAt `gorm:"index" faker:"-"`
}

func (m *User) TableName() string {
	return UserTableName
}

func (m *User) ToEntity() *entity.User {
	return &entity.User{
		UUID:                  m.UUID,
		TraceID:               m.TraceID,
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

func NewUser(e entity.SaveUser) *User {
	return &User{
		Username:              e.Username,
		TraceID:               e.TraceID,
		Password:              e.Password,
		Email:                 e.Email,
		PhoneNumber:           e.PhoneNumber,
		FirstName:             e.FirstName,
		LastName:              e.LastName,
		IsEmailVerified:       false,
		IsPhoneNumberVerified: false,
	}
}
