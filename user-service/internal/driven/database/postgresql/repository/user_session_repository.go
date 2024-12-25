package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database/postgresql"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/fx"
	"gorm.io/gorm/clause"
)

type UserSessionRepositoryParam struct {
	fx.In

	Database *database.Database
	Tracer   tracer.Tracer
	Logger   logger.Logger
	Util     *util.Util
}

type UserSessionRepositoryResult struct {
	fx.Out

	UserSessionStorage shared_service.UserSessionStorage
}

type userSessionRepositoryImpl struct {
	database *database.Database
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewUserSessionRepository(p UserSessionRepositoryParam) UserSessionRepositoryResult {
	return UserSessionRepositoryResult{
		UserSessionStorage: &userSessionRepositoryImpl{
			database: p.Database,
			tracer:   p.Tracer,
			logger:   p.Logger,
			util:     p.Util,
		},
	}
}

func (r *userSessionRepositoryImpl) WithTx(tx *database.Transaction) shared_service.UserSessionStorage {
	if tx == nil {
		return r
	}

	return &userSessionRepositoryImpl{
		database: &database.Database{
			DB:     tx.DB,
			Logger: r.logger,
		},
		tracer: r.tracer,
		logger: r.logger,
		util:   r.util,
	}
}

func (r *userSessionRepositoryImpl) SaveSession(ctx context.Context, createModel entity.SaveUserSession) (*entity.UserSession, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	createMap, err := r.util.StructUtil.ConvertSqlStructToMap(createModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	var userSessionModel *model.UserSession

	err = r.database.DB.WithContext(ctx).Model(&userSessionModel).Create(createMap).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == postgresql.UniqueViolation {
				if strings.Contains(err.Error(), "email") {
					return nil, auth_service.ErrDuplicateEmail
				} else if strings.Contains(err.Error(), "phone_number") {
					return nil, auth_service.ErrDuplicatePhoneNumber
				}
			}
		}
	}

	return userSessionModel.ToEntity(), err
}

func (r *userSessionRepositoryImpl) FindSession(ctx context.Context, findModel entity.FindUserSession) (*entity.UserSession, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	var userSessionModel model.UserSession
	err = r.database.DB.WithContext(ctx).Where(findMap).First(&userSessionModel).Error
	if err != nil {
		return nil, err
	}

	return userSessionModel.ToEntity(), err
}

func (r *userSessionRepositoryImpl) RevokeSession(ctx context.Context, refreshToken string) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	err := r.database.DB.
		WithContext(ctx).
		Model(&model.UserSession{}).
		Where(model.UserSession{RefreshToken: refreshToken}).
		Clauses(clause.Returning{}).
		Updates(model.UserSession{IsRevoked: true}).
		Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
	}

	return err
}

func (r *userSessionRepositoryImpl) RevokeManySession(ctx context.Context, refreshTokens []string) error {
	return r.database.DB.
		WithContext(ctx).
		Where("refresh_token IN ?", refreshTokens).
		Delete(&model.UserSession{}).
		Error
}
