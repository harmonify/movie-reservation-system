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
		FindUser(ctx context.Context, findModel entity.FindUser) (*entity.User, error)
		// FindUserWithResult is a generic function to find a user with a result model
		// The result model parameter should be a pointer to a struct
		FindUserWithResult(ctx context.Context, findModel entity.FindUser, resultModel interface{}) error
		SaveUser(ctx context.Context, createModel entity.SaveUser) (*entity.User, error)
		UpdateUser(ctx context.Context, findModel entity.FindUser, updateModel entity.UpdateUser) (*entity.User, error)
		SoftDeleteUser(ctx context.Context, findModel entity.FindUser) error
	}

	UserSessionStorage interface {
		WithTx(tx *database.Transaction) UserSessionStorage
		FindSession(ctx context.Context, findModel entity.FindUserSession) (*entity.UserSession, error)
		SaveSession(ctx context.Context, createModel entity.SaveUserSession) (*entity.UserSession, error)
		RevokeSession(ctx context.Context, refreshToken string) (err error)
		RevokeManySession(ctx context.Context, refreshTokens []string) (err error)
		SoftDeleteSession(ctx context.Context, findModel entity.FindUserSession) error
	}

	UserKeyStorage interface {
		WithTx(tx *database.Transaction) UserKeyStorage
		FindUserKey(ctx context.Context, findModel entity.FindUserKey) (*entity.UserKey, error)
		SaveUserKey(ctx context.Context, createModel entity.SaveUserKey) (*entity.UserKey, error)
		UpdateUserKey(ctx context.Context, findModel entity.FindUserKey, updateModel entity.UpdateUserKey) (*entity.UserKey, error)
		SoftDeleteUserKey(ctx context.Context, findModel entity.FindUserKey) error
	}

	OutboxStorage interface {
		WithTx(tx *database.Transaction) OutboxStorage
		SaveOutbox(ctx context.Context, createModel entity.SaveUserOutbox) (*entity.UserOutbox, error)
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
)
