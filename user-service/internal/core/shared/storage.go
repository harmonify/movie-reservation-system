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
	}

	UserKeyStorage interface {
		WithTx(tx *database.Transaction) UserKeyStorage
		FindUserKey(ctx context.Context, findModel entity.FindUserKey) (*entity.UserKey, error)
		SaveUserKey(ctx context.Context, createModel entity.SaveUserKey) (*entity.UserKey, error)
		UpdateUserKey(ctx context.Context, findModel entity.FindUserKey, updateModel entity.UpdateUserKey) (*entity.UserKey, error)
		SoftDeleteUserKey(ctx context.Context, findModel entity.FindUserKey) error
	}

	OtpStorage interface {
		SaveEmailVerificationToken(ctx context.Context, p SaveEmailVerificationTokenParam) error
		GetEmailVerificationToken(ctx context.Context, email string) (string, error)
		DeleteEmailVerificationToken(ctx context.Context, email string) (bool, error)
		SavePhoneOtp(ctx context.Context, p SavePhoneOtpParam) error
		GetPhoneOtp(ctx context.Context, phoneNumber string) (string, error)
		DeletePhoneOtp(ctx context.Context, phoneNumber string) (bool, error)
		IncrementPhoneOtpAttempt(ctx context.Context, phoneNumber string) error
		GetPhoneOtpAttempt(ctx context.Context, phoneNumber string) (int, error)
		DeletePhoneOtpAttempt(ctx context.Context, phoneNumber string) (bool, error)
	}

	OutboxStorage interface {
		WithTx(tx *database.Transaction) OutboxStorage
		SaveOutbox(ctx context.Context, createModel entity.SaveUserOutbox) (*entity.UserOutbox, error)
	}
)
