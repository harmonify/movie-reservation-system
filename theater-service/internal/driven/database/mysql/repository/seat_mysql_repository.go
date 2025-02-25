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

type seatRepositoryImpl struct {
	database *database.Database
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewSeatRepository(
	database *database.Database,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared.SeatStorage {
	return &seatRepositoryImpl{
		database: database,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *seatRepositoryImpl) WithTx(tx *database.Transaction) shared.SeatStorage {
	if tx == nil {
		return r
	}
	return NewSeatRepository(
		r.database.WithTx(tx),
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *seatRepositoryImpl) SaveSeat(ctx context.Context, create *entity.SaveSeat) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	result := r.database.DB.
		WithContext(ctx).
		Create(&create)

	return result.Error
}

func (r *seatRepositoryImpl) FindSeat(ctx context.Context, find *entity.FindSeat) (*entity.Seat, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return nil, err
	}

	seat := &entity.Seat{}
	result := r.database.DB.WithContext(ctx).Where(findMap).First(&seat)
	err = result.Error
	if err != nil {
		return nil, err
	}

	return seat, err
}

func (r *seatRepositoryImpl) UpdateSeat(ctx context.Context, find *entity.FindSeat, update *entity.UpdateSeat) error {
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
		Model(&entity.Seat{}).
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

func (r *seatRepositoryImpl) SoftDeleteSeat(ctx context.Context, find *entity.FindSeat) error {
	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(findMap).
		Delete(&entity.Seat{})

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

func (r *seatRepositoryImpl) CountRoomSeats(ctx context.Context, roomIds []string) ([]*entity.CountRoomSeats, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var countRoomSeats []*entity.CountRoomSeats
	result := r.database.DB.
		WithContext(ctx).
		Model(&entity.Seat{}).
		Select("room_id, count(*) as count").
		Where("room_id IN ?", roomIds).
		Group("room_id").
		Find(&countRoomSeats)

	err := result.Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return countRoomSeats, nil
}

func (r *seatRepositoryImpl) FindShowtimeAvailableSeats(ctx context.Context, findModel *entity.FindShowtimeAvailableSeats) ([]*entity.Seat, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var seats []*entity.Seat
	result := r.database.DB.
		WithContext(ctx).
		Model(&entity.Seat{}).
		Where("showtime_id = ?", findModel.ShowtimeID).
		Where(
			"seat_id NOT IN (?)",
			r.database.DB.
				WithContext(ctx).
				Model(&entity.Ticket{}).
				Select("DISTINCT seat_id").
				Where("showtime_id = ?", findModel.ShowtimeID),
		).
		Find(&seats)

	err := result.Error
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return nil, err
	}

	return seats, nil
}
