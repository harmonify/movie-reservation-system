package shared

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
)

type (
	Storage interface {
		Transaction(fc func(tx *database.Transaction) error) error
	}

	UserStorage interface {
		WithTx(tx *database.Transaction) UserStorage
		GetUser(ctx context.Context, getModel entity.GetUser) (*entity.User, error)
		// GetUserWithResult is a generic function to Get a user with a result model
		// The result model parameter should be a pointer to a struct
		GetUserWithResult(ctx context.Context, getModel entity.GetUser, resultModel interface{}) error
		SaveUser(ctx context.Context, createModel entity.SaveUser) (*entity.User, error)
		UpdateUser(ctx context.Context, getModel entity.GetUser, updateModel entity.UpdateUser) (*entity.User, error)
		SoftDeleteUser(ctx context.Context, getModel entity.GetUser) error
	}

	UserSessionStorage interface {
		WithTx(tx *database.Transaction) UserSessionStorage
		GetSession(ctx context.Context, getModel entity.GetUserSession) (*entity.UserSession, error)
		SaveSession(ctx context.Context, createModel entity.SaveUserSession) (*entity.UserSession, error)
		RevokeSession(ctx context.Context, refreshToken string) (err error)
		RevokeManySession(ctx context.Context, refreshTokens []string) (err error)
		SoftDeleteSession(ctx context.Context, getModel entity.GetUserSession) error
	}

	UserKeyStorage interface {
		WithTx(tx *database.Transaction) UserKeyStorage
		GetUserKey(ctx context.Context, getModel entity.GetUserKey) (*entity.UserKey, error)
		SaveUserKey(ctx context.Context, createModel entity.SaveUserKey) (*entity.UserKey, error)
		UpdateUserKey(ctx context.Context, getModel entity.GetUserKey, updateModel entity.UpdateUserKey) (*entity.UserKey, error)
		SoftDeleteUserKey(ctx context.Context, getModel entity.GetUserKey) error
	}

	OutboxStorage interface {
		WithTx(tx *database.Transaction) OutboxStorage
		SaveOutbox(ctx context.Context, createModel entity.SaveUserOutbox) (*entity.UserOutbox, error)
	}

	// RoleStorage interface {
	// 	WithTx(tx *database.Transaction) RoleStorage
	// 	GetRole(ctx context.Context, getModel entity.GetRole) (*entity.Role, error)
	// 	SaveRole(ctx context.Context, createModel entity.SaveRole) (*entity.Role, error)
	// 	UpdateRole(ctx context.Context, getModel entity.GetRole, updateModel entity.UpdateRole) (*entity.Role, error)
	// 	SoftDeleteRole(ctx context.Context, getModel entity.GetRole) error
	// }

	UserRoleStorage interface {
		WithTx(tx *database.Transaction) UserRoleStorage
		SearchUserRoles(ctx context.Context, searchModel entity.SearchUserRoles) ([]*entity.UserRole, error)
		SaveUserRoles(ctx context.Context, createModel entity.SaveUserRoles) ([]*entity.UserRole, error)
		SoftDeleteUserRoles(ctx context.Context, searchModel entity.SearchUserRoles) error
	}

	OtpCache interface {
		SaveEmailVerificationCode(ctx context.Context, p SaveEmailVerificationCodeParam) error
		GetEmailVerificationCode(ctx context.Context, uuid string) (string, error)
		DeleteEmailVerificationCode(ctx context.Context, uuid string) (bool, error)
		SavePhoneNumberVerificationOtp(ctx context.Context, p SavePhoneNumberVerificationOtpParam) error
		GetPhoneNumberVerificationOtp(ctx context.Context, uuid string) (string, error)
		DeletePhoneNumberVerificationOtp(ctx context.Context, uuid string) (bool, error)
		IncrementPhoneNumberVerificationAttempt(ctx context.Context, uuid string) error
		GetPhoneNumberVerificationAttempt(ctx context.Context, uuid string) (int, error)
		DeletePhoneNumberVerificationAttempt(ctx context.Context, uuid string) (bool, error)
	}

	OtpCacheV2 interface {
		// SaveOtp saves the otp code for the given uuid and otp type
		SaveOtp(ctx context.Context, uuid string, otpType OtpType, code string) error
		// GetOtp Gets the otp for the given uuid and otp type
		// If the otp is not found, it will return otp_service.OtpNotFoundError
		GetOtp(ctx context.Context, uuid string, otpType OtpType) (*Otp, error)
		// DeleteOtp deletes the otp for the given uuid and otp type
		DeleteOtp(ctx context.Context, uuid string, otpType OtpType) (bool, error)
		// IncrementOtpAttempt increments the attempt count for the given uuid and otp type
		// If the otp is not found, it will return otp_service.OtpNotFoundError
		IncrementOtpAttempt(ctx context.Context, uuid string, otpType OtpType) (*Otp, error)
	}
)
