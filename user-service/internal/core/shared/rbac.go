package shared

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/database"
)

var (
	ActionGet    Action = "GET"
	ActionPost   Action = "POST"
	ActionPut    Action = "PUT"
	ActionDelete Action = "DELETE"

	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type (
	Action string
	Role   string

	RbacStorage interface {
		WithTx(tx *database.Transaction) RbacStorage
		CheckPermission(ctx context.Context, p CheckPermissionParam) (bool, error)
		GrantPermission(ctx context.Context, p GrantPermissionParam) (bool, error)
		BulkGrantPermission(ctx context.Context, p []GrantPermissionParam) (bool, error)
		RevokePermission(ctx context.Context, p RevokePermissionParam) (bool, error)
		GrantRole(ctx context.Context, p GrantRoleParam) (bool, error)
	}

	CheckPermissionParam struct {
		UUID     string
		Domain   string
		Resource string
		Action   Action
	}

	GrantPermissionParam struct {
		UUID     string
		Domain   string
		Resource string
		Action   Action
	}

	RevokePermissionParam struct {
		UUID     string
		Domain   string
		Resource string
		Action   Action
	}

	GrantRoleParam struct {
		UUID string
		Role Role
	}
)
