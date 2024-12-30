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
) shared_service.UserKeyStorage {
	return &userKeyRepositoryImpl{
		database: database,
		pgErrTl:  pgErrTl,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *userKeyRepositoryImpl) WithTx(tx *database.Transaction) shared_service.UserKeyStorage {
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

	userKeyModel := (&model.UserKey{}).FromSaveEntity(createModel)

	result := r.database.DB.
		WithContext(ctx).
		Create(userKeyModel)
	err := r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	return userKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) FindUserKey(ctx context.Context, findModel entity.FindUserKey) (*entity.UserKey, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	userKeyModel := model.UserKey{}
	result := r.database.DB.WithContext(ctx).Where(findMap).First(&userKeyModel)
	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	return userKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) UpdateUserKey(ctx context.Context, findModel entity.FindUserKey, updateModel entity.UpdateUserKey) (*entity.UserKey, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	updateMap, err := r.util.StructUtil.ConvertSqlStructToMap(updateModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	userKeyModel := model.UserKey{}
	result := r.database.DB.
		WithContext(ctx).
		Model(&userKeyModel).
		Where(findMap).
		Clauses(clause.Returning{}).
		Updates(updateMap)

	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error())
		return nil, err
	}

	return userKeyModel.ToEntity(), err
}

func (r *userKeyRepositoryImpl) SoftDeleteUserKey(ctx context.Context, findModel entity.FindUserKey) error {
	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(findModel)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(findMap).
		Delete(&model.UserKey{})

	err = r.pgErrTl.Translate(result.Error)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error())
		return err
	}

	rowsAffected := result.RowsAffected
	if rowsAffected <= 0 {
		err := database.NewRecordNotFoundError(gorm.ErrRecordNotFound)
		r.logger.WithCtx(ctx).Error(err.Error())
		return err
	}

	return nil
}
