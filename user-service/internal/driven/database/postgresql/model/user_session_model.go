package model

import (
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"gorm.io/gorm"
)

type UserSession struct {
	gorm.Model

	UserUUID     string `gorm:"index"`
	TraceID      string `gorm:"uniqueIndex:idx_user_trace_id;unique"`
	RefreshToken string // hashed
	IsRevoked    bool   `gorm:"default:false"`
	ExpiredAt    time.Time
	IpAddress    sql.NullString
	UserAgent    sql.NullString

	User User `gorm:"foreignKey:UserUUID"`
}

func (m *UserSession) TableName() string {
	return "user_sessions"
}

func (m *UserSession) ToEntity() *entity.UserSession {
	return &entity.UserSession{
		ID:           uint64(m.ID),
		UserUUID:     m.UserUUID,
		RefreshToken: m.RefreshToken,
		IsRevoked:    m.IsRevoked,
		ExpiredAt:    m.ExpiredAt,
		IpAddress:    m.IpAddress,
		UserAgent:    m.UserAgent,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    sql.NullTime(m.DeletedAt),
	}
}

func NewUserSession(e entity.SaveUserSession) *UserSession {
	return &UserSession{
		UserUUID:     e.UserUUID,
		TraceID:      e.TraceID,
		RefreshToken: e.RefreshToken,
		IsRevoked:    false,
		ExpiredAt:    e.ExpiredAt,
		IpAddress:    e.IpAddress,
		UserAgent:    e.UserAgent,
	}
}
