package repository

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userKeyRepositoryImpl struct {
	database *database.Database
	pgErrTl  database.PostgresqlErrorTranslator
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewUserKeyRepository(
	database *database.Database,
	pgErrTl database.PostgresqlErrorTranslator,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared.UserKeyStorage {
	return &userKeyRepositoryImpl{
		database: database,
		pgErrTl:  pgErrTl,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *userKeyRepositoryImpl) WithTx(tx *database.Transaction) shared.UserKeyStorage {
	if tx == nil {
		return r
	}
	return NewUserKeyRepository(
		r.database.WithTx(tx),
		r.pgErrTl,
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *userKeyRepositoryImpl) SaveUserKey(ctx context.Context, createModel entity.SaveUserKey) (*entity.UserKey, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	userKeyModel := model.NewUserKey(createModel)

	result := r.database.DB.
		WithContext(ctx).
		Create(userKeyModel)
	err := r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return userKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) GetUserKey(ctx context.Context, getModel entity.GetUserKey) (*entity.UserKey, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	getMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, getModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	userKeyModel := model.UserKey{}
	result := r.database.DB.WithContext(ctx).Where(getMap).First(&userKeyModel)
	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return userKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) UpdateUserKey(ctx context.Context, getModel entity.GetUserKey, updateModel entity.UpdateUserKey) (*entity.UserKey, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	getMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, getModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	updateMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, updateModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	userKeyModel := model.UserKey{}
	result := r.database.DB.
		WithContext(ctx).
		Model(&userKeyModel).
		Where(getMap).
		Clauses(clause.Returning{}).
		Updates(updateMap)

	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return userKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) SoftDeleteUserKey(ctx context.Context, getModel entity.GetUserKey) error {
	getMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, getModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(getMap).
		Delete(&model.UserKey{})

	err = r.pgErrTl.Translate(result.Error)
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
