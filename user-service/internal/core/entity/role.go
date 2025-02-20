package entity

import (
	"database/sql"
	"time"
)

// FixedRole is the defined roles in the system
// Its value represents the role ID in the storage
type FixedRole uint

func (r FixedRole) Value() uint {
	return uint(r)
}

// String returns the string representation of the role
func (r FixedRole) String() string {
	return [...]string{"", "admin", "user"}[r]
}

const (
	RoleAdmin FixedRole = iota + 1
	RoleUser
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
