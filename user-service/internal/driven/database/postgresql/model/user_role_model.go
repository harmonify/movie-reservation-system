package model

import (
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"gorm.io/gorm"
)

type UserRole struct {
	ID        uint           `json:"id"`
	UserUUID  string         `json:"user_id"`
	RoleID    uint           `json:"role_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func (m *UserRole) TableName() string {
	return "user_roles"
}

func (m *UserRole) ToEntity() *entity.UserRole {
	return &entity.UserRole{
		ID:        m.ID,
		UserUUID:  m.UserUUID,
		RoleID:    m.RoleID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: sql.NullTime(m.DeletedAt),
	}
}

func NewUserRole(userUUID string, roleId uint) *UserRole {
	return &UserRole{
		UserUUID: userUUID,
		RoleID:   roleId,
	}
}
