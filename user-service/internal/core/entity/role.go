package entity

import (
	"database/sql"
	"time"
)

type RoleName string

const (
	RoleAdmin RoleName = "admin"
	RoleUser  RoleName = "user"
)

type Role struct {
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type GetRole struct {
	ID uint
}

type SaveRole struct {
	Name string
}

type UpdateRole struct {
	Name string
}
