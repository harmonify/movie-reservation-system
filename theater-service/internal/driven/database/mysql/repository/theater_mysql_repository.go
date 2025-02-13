package repository

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/shared"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type theaterRepositoryImpl struct {
	database *database.Database
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewTheaterRepository(
	database *database.Database,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared.TheaterStorage {
	return &theaterRepositoryImpl{
		database: database,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *theaterRepositoryImpl) WithTx(tx *database.Transaction) shared.TheaterStorage {
	if tx == nil {
		return r
	}
	return NewTheaterRepository(
		r.database.WithTx(tx),
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *theaterRepositoryImpl) SaveTheater(ctx context.Context, create *entity.SaveTheater) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	result := r.database.DB.
		WithContext(ctx).
		Create(&create)

	return result.Error
}

func (r *theaterRepositoryImpl) UpdateTheater(ctx context.Context, find *entity.FindOneTheater, update *entity.UpdateTheater) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	updateMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, update)
	if err != nil {
		return err
	}

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Model(&entity.Theater{}).
		Where(findMap).
		Updates(updateMap)

	err = result.Error
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

func (r *theaterRepositoryImpl) SoftDeleteTheater(ctx context.Context, find *entity.FindOneTheater) error {
	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(findMap).
		Delete(&entity.Theater{})

	err = result.Error
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

func (r *theaterRepositoryImpl) FindOneTheater(ctx context.Context, find *entity.FindOneTheater) (*entity.Theater, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return nil, err
	}

	theater := &entity.Theater{}
	result := r.database.DB.WithContext(ctx).Where(findMap).First(&theater)
	err = result.Error
	if err != nil {
		return nil, err
	}

	return theater, err
}

func (r *theaterRepositoryImpl) FindManyTheaters(ctx context.Context, find *entity.FindManyTheaters) ([]*entity.Theater, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return nil, err
	}

	theaters := []*entity.Theater{}
	result := r.database.DB.WithContext(ctx).Where(findMap).Find(&theaters)
	err = result.Error
	if err != nil {
		return nil, err
	}

	return theaters, err
}
