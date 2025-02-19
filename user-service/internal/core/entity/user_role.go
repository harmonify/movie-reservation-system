package entity

import (
	"time"
)

type UserRole struct {
	ID        uint      `json:"id"`
	UserUUID  string    `json:"user_id"`
	RoleID    uint      `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserRole struct {
	ID       uint
	UserUUID string
	RoleID   uint
}

type SaveUserRole struct {
	UserUUID string
	RoleName RoleName
}
