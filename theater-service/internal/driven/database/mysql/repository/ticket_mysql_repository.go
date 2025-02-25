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

type ticketRepositoryImpl struct {
	database *database.Database
	tracer   tracer.Tracer
	logger   logger.Logger
	util     *util.Util
}

func NewTicketRepository(
	database *database.Database,
	tracer tracer.Tracer,
	logger logger.Logger,
	util *util.Util,
) shared.TicketStorage {
	return &ticketRepositoryImpl{
		database: database,
		tracer:   tracer,
		logger:   logger,
		util:     util,
	}
}

func (r *ticketRepositoryImpl) WithTx(tx *database.Transaction) shared.TicketStorage {
	if tx == nil {
		return r
	}
	return NewTicketRepository(
		r.database.WithTx(tx),
		r.tracer,
		r.logger,
		r.util,
	)
}

func (r *ticketRepositoryImpl) SaveTicket(ctx context.Context, create *entity.SaveTicket) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	result := r.database.DB.
		WithContext(ctx).
		Create(&create)

	return result.Error
}

func (r *ticketRepositoryImpl) UpdateTicket(ctx context.Context, find *entity.FindOneTicket, update *entity.UpdateTicket) error {
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
		Model(&entity.Ticket{}).
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

func (r *ticketRepositoryImpl) SoftDeleteTicket(ctx context.Context, find *entity.FindOneTicket) error {
	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		r.logger.WithCtx(ctx).Error(err.Error(), zap.Error(err))
		return err
	}

	result := r.database.DB.
		WithContext(ctx).
		Where(findMap).
		Delete(&entity.Ticket{})

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

func (r *ticketRepositoryImpl) FindOneTicket(ctx context.Context, find *entity.FindOneTicket) (*entity.Ticket, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return nil, err
	}

	ticket := &entity.Ticket{}
	result := r.database.DB.WithContext(ctx).Where(findMap).First(&ticket)
	err = result.Error
	if err != nil {
		return nil, err
	}

	return ticket, err
}

func (r *ticketRepositoryImpl) FindManyTickets(ctx context.Context, find *entity.FindManyTickets) ([]*entity.Ticket, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	findMap, err := r.util.StructUtil.ConvertSqlStructToMap(ctx, find)
	if err != nil {
		return nil, err
	}

	tickets := []*entity.Ticket{}
	result := r.database.DB.WithContext(ctx).Where(findMap).Find(&tickets)
	err = result.Error
	if err != nil {
		return nil, err
	}

	return tickets, err
}

func (r *ticketRepositoryImpl) CountShowtimeTickets(ctx context.Context, showtimeIds []string) ([]*entity.CountShowtimeTicket, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	tickets := make([]*entity.CountShowtimeTicket, 0)
	result := r.database.DB.WithContext(ctx).
		Model(&entity.Ticket{}).
		Select("showtime_id, count(*) as count").
		Where("showtime_id IN (?)", showtimeIds).
		Group("showtime_id").
		Find(&tickets)

	err := result.Error
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
