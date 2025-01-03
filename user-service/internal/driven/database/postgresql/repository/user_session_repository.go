package repository

import (
	"context"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSessionRepositoryImpl struct {
	database *database.Database
	pgErrTl  database.PostgresqlErrorTranslator
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewUserSessionRepository(
	database *database.Database,
	pgErrTl database.PostgresqlErrorTranslator,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared_service.UserSessionStorage {
	return &userSessionRepositoryImpl{
		database: database,
		pgErrTl:  pgErrTl,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *userSessionRepositoryImpl) WithTx(tx *database.Transaction) shared_service.UserSessionStorage {
	if tx == nil {
		return r
	}
	return NewUserSessionRepository(
		r.database.WithTx(tx),
		r.pgErrTl,
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *userSessionRepositoryImpl) SaveSession(ctx context.Context, createModel entity.SaveUserSession) (*entity.UserSession, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	userSessionModel := (&model.UserSession{}).FromSaveEntity(createModel)

	result := r.database.DB.
		WithContext(ctx).
		Create(userSessionModel)
	err := r.pgErrTl.Translate(result.Error)
	if err != nil {
		return nil, err
	}

	return userSessionModel.ToEntity(), nil
}

func (r *userSessionRepositoryImpl) FindSession(ctx context.Context, findModel entity.FindUserSession) (*entity.UserSession, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	userSessionModel := model.UserSession{}
	result := r.database.DB.WithContext(ctx).Where(findMap).First(&userSessionModel)
	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		return nil, err
	}

	return userSessionModel.ToEntity(), err
}

func (r *userSessionRepositoryImpl) RevokeSession(ctx context.Context, refreshToken string) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	result := r.database.DB.
		WithContext(ctx).
		Model(&model.UserSession{}).
		Where(model.UserSession{RefreshToken: refreshToken}).
		Clauses(clause.Returning{}).
		Updates(model.UserSession{IsRevoked: true})
	err := r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	return err
}

func (r *userSessionRepositoryImpl) RevokeManySession(ctx context.Context, refreshTokens []string) error {
	result := r.database.DB.
		WithContext(ctx).
		Where("refresh_token IN ?", refreshTokens).
		Updates(model.UserSession{IsRevoked: true})

	err := r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	return nil
}
