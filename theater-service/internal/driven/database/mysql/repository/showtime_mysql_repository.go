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

type showtimeRepositoryImpl struct {
	database *database.Database
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewShowtimeRepository(
	database *database.Database,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared.ShowtimeStorage {
	return &showtimeRepositoryImpl{
		database: database,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *showtimeRepositoryImpl) WithTx(tx *database.Transaction) shared.ShowtimeStorage {
	if tx == nil {
		return r
	}
	return NewShowtimeRepository(
		r.database.WithTx(tx),
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *showtimeRepositoryImpl) SaveShowtime(ctx context.Context, create *entity.SaveShowtime) (*entity.SaveShowtimeResult, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	showtime := entity.NewShowtime(create)

	err := r.database.DB.
		WithContext(ctx).
		Create(&showtime).
		Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return &entity.SaveShowtimeResult{
		ShowtimeID: showtime.ShowtimeID,
	}, nil
}

func (r *showtimeRepositoryImpl) UpdateShowtime(ctx context.Context, find *entity.FindOneShowtime, update *entity.UpdateShowtime) error {
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
		Model(&entity.Showtime{}).
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

func (r *showtimeRepositoryImpl) SoftDeleteShowtime(ctx context.Context, find *entity.FindOneShowtime) error {
	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(findMap).
		Delete(&entity.Showtime{})

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

func (r *showtimeRepositoryImpl) FindOneShowtime(ctx context.Context, find *entity.FindOneShowtime) (*entity.Showtime, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return nil, err
	}

	var showtime entity.Showtime
	result := r.database.DB.
		WithContext(ctx).
		Where(findMap).
		First(&showtime)

	err = result.Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return &showtime, nil
}

func (r *showtimeRepositoryImpl) FindManyShowtimes(ctx context.Context, find *entity.FindManyShowtimes) (*entity.FindManyShowtimesResult, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return nil, err
	}
	delete(findMap, "start_time_gte")
	delete(findMap, "start_time_lte")

	query := r.database.DB.WithContext(ctx).Where(findMap)

	if find.StartTimeGte.Valid {
		query = query.Where("start_time >= ?", find.StartTimeGte.Time)
	}
	if find.StartTimeLte.Valid {
		query = query.Where("start_time <= ?", find.StartTimeLte.Time)
	}

	var showtimes []*entity.Showtime
	result := query.Find(&showtimes)

	err = result.Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	var totalResults int64
	result = query.Count(&totalResults)
	if result.Error != nil {
		r.logger.WithCtx(ctx).Error(result.Error.Error(), zap.Error(result.Error))
		return nil, result.Error
	}

	return &entity.FindManyShowtimesResult{
		Showtimes: showtimes,
	}, nil
}
