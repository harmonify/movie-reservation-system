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
	AdminTheaterService interface {
		SearchTheaters(ctx context.Context, findModel *entity.FindManyTheaters) (*entity.FindManyTheatersResult, error)
		GetTheaterByID(ctx context.Context, findModel *entity.FindOneTheater) (*entity.Theater, error)
		SaveTheater(ctx context.Context, saveModel *entity.SaveTheater) (*entity.SaveTheaterResult, error)
		UpdateTheater(ctx context.Context, findModel *entity.FindOneTheater, updateModel *entity.UpdateTheater) error
		SoftDeleteTheater(ctx context.Context, findModel *entity.FindOneTheater) error
	}

	AdminTheaterServiceParam struct {
		fx.In
		Logger         logger.Logger
		Tracer         tracer.Tracer
		TheaterStorage shared.TheaterStorage
	}

	AdminTheaterServiceResult struct {
		fx.Out

		AdminTheaterService AdminTheaterService
	}

	adminTheaterServiceImpl struct {
		logger         logger.Logger
		tracer         tracer.Tracer
		theaterStorage shared.TheaterStorage
	}
)

func NewAdminTheaterService(p AdminTheaterServiceParam) AdminTheaterServiceResult {
	s := &adminTheaterServiceImpl{
		logger:         p.Logger,
		tracer:         p.Tracer,
		theaterStorage: p.TheaterStorage,
	}

	return AdminTheaterServiceResult{
		AdminTheaterService: s,
	}
}

func (s *adminTheaterServiceImpl) SearchTheaters(ctx context.Context, p *entity.FindManyTheaters) (*entity.FindManyTheatersResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.theaterStorage.FindManyTheaters(ctx, p)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to search theaters", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *adminTheaterServiceImpl) GetTheaterByID(ctx context.Context, findModel *entity.FindOneTheater) (*entity.Theater, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.theaterStorage.FindOneTheater(ctx, findModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get theaters", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *adminTheaterServiceImpl) SaveTheater(ctx context.Context, saveModel *entity.SaveTheater) (*entity.SaveTheaterResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.theaterStorage.SaveTheater(ctx, saveModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to save theaters", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (s *adminTheaterServiceImpl) UpdateTheater(ctx context.Context, findModel *entity.FindOneTheater, updateModel *entity.UpdateTheater) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	err := s.theaterStorage.UpdateTheater(ctx, findModel, updateModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to update theaters", zap.Error(err))
		return err
	}

	return nil
}

func (s *adminTheaterServiceImpl) SoftDeleteTheater(ctx context.Context, findModel *entity.FindOneTheater) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	err := s.theaterStorage.SoftDeleteTheater(ctx, findModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to soft delete theaters", zap.Error(err))
		return err
	}

	return nil
}
