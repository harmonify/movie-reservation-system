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
	TheaterService interface {
		SearchTheaters(ctx context.Context, findModel *entity.FindManyTheaters) (*entity.FindManyTheatersResult, error)
	}

	TheaterServiceParam struct {
		fx.In
		Logger         logger.Logger
		Tracer         tracer.Tracer
		TheaterStorage shared.TheaterStorage
	}

	TheaterServiceResult struct {
		fx.Out

		TheaterService TheaterService
	}

	TheaterServiceImpl struct {
		logger         logger.Logger
		tracer         tracer.Tracer
		theaterStorage shared.TheaterStorage
	}
)

func NewTheaterService(p TheaterServiceParam) TheaterServiceResult {
	s := &TheaterServiceImpl{
		logger:         p.Logger,
		tracer:         p.Tracer,
		theaterStorage: p.TheaterStorage,
	}

	return TheaterServiceResult{
		TheaterService: s,
	}
}

func (s *TheaterServiceImpl) SearchTheaters(ctx context.Context, p *entity.FindManyTheaters) (*entity.FindManyTheatersResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.theaterStorage.FindManyTheaters(ctx, p)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to search theaters", zap.Error(err))
		return nil, err
	}

	return res, nil
}
