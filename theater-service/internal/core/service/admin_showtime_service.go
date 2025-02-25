package service

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/shared"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	AdminShowtimeService interface {
		SearchShowtimes(ctx context.Context, findModel *entity.FindManyShowtimes) (*entity.FindManyShowtimesResult, error)
		GetShowtimeByID(ctx context.Context, findModel *entity.FindOneShowtime) (*entity.Showtime, error)
		SaveShowtime(ctx context.Context, saveModel *entity.SaveShowtime) (*entity.SaveShowtimeResult, error)
		UpdateShowtime(ctx context.Context, findModel *entity.FindOneShowtime, updateModel *entity.UpdateShowtime) error
		SoftDeleteShowtime(ctx context.Context, findModel *entity.FindOneShowtime) error
	}

	AdminShowtimeServiceParam struct {
		fx.In
		Logger          logger.Logger
		Tracer          tracer.Tracer
		ShowtimeStorage shared.ShowtimeStorage
		SeatStorage     shared.SeatStorage
		TicketStorage   shared.TicketStorage
	}

	AdminShowtimeServiceResult struct {
		fx.Out

		AdminShowtimeService AdminShowtimeService
	}

	adminShowtimeServiceImpl struct {
		logger          logger.Logger
		tracer          tracer.Tracer
		showtimeStorage shared.ShowtimeStorage
	}
)

func NewAdminShowtimeService(p AdminShowtimeServiceParam) AdminShowtimeServiceResult {
	s := &adminShowtimeServiceImpl{
		logger:          p.Logger,
		tracer:          p.Tracer,
		showtimeStorage: p.ShowtimeStorage,
	}

	return AdminShowtimeServiceResult{
		AdminShowtimeService: s,
	}
}

func (s *adminShowtimeServiceImpl) SearchShowtimes(ctx context.Context, p *entity.FindManyShowtimes) (*entity.FindManyShowtimesResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.showtimeStorage.FindManyShowtimes(ctx, p)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to search showtimes", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *adminShowtimeServiceImpl) GetShowtimeByID(ctx context.Context, findModel *entity.FindOneShowtime) (*entity.Showtime, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.showtimeStorage.FindOneShowtime(ctx, findModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get showtimes", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *adminShowtimeServiceImpl) SaveShowtime(ctx context.Context, saveModel *entity.SaveShowtime) (*entity.SaveShowtimeResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.showtimeStorage.SaveShowtime(ctx, saveModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save showtimes", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *adminShowtimeServiceImpl) UpdateShowtime(ctx context.Context, findModel *entity.FindOneShowtime, updateModel *entity.UpdateShowtime) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	err := s.showtimeStorage.UpdateShowtime(ctx, findModel, updateModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to update showtimes", zap.Error(err))
		return err
	}

	return nil
}

func (s *adminShowtimeServiceImpl) SoftDeleteShowtime(ctx context.Context, findModel *entity.FindOneShowtime) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	err := s.showtimeStorage.SoftDeleteShowtime(ctx, findModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to soft delete showtimes", zap.Error(err))
		return err
	}

	return nil
}
