package entity

import (
	"database/sql"
	"time"
)

type UserRole struct {
	ID        uint         `json:"id"`
	UserUUID  string       `json:"user_id"`
	RoleID    uint         `json:"role_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty"`
}

type SearchUserRoles struct {
	ID       sql.NullInt32
	UserUUID sql.NullString
	RoleID   sql.NullInt32
}

type SaveUserRoles struct {
	UserUUID string
	RoleID   []uint
}
