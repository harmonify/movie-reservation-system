package service

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/shared"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	SeatService interface {
		GetAvailableSeats(ctx context.Context, req *theater_proto.GetAvailableSeatsRequest) (*theater_proto.GetAvailableSeatsResponse, error)
	}

	SeatServiceParam struct {
		fx.In
		Logger          logger.Logger
		Tracer          tracer.Tracer
		TheaterStorage  shared.TheaterStorage
		ShowtimeStorage shared.ShowtimeStorage
		SeatStorage     shared.SeatStorage
		TicketStorage   shared.TicketStorage
	}

	SeatServiceResult struct {
		fx.Out

		SeatService SeatService
	}

	SeatServiceImpl struct {
		logger          logger.Logger
		tracer          tracer.Tracer
		theaterStorage  shared.TheaterStorage
		showtimeStorage shared.ShowtimeStorage
		seatStorage     shared.SeatStorage
		ticketStorage   shared.TicketStorage
	}
)

func NewSeatService(p SeatServiceParam) SeatServiceResult {
	s := &SeatServiceImpl{
		logger:          p.Logger,
		tracer:          p.Tracer,
		theaterStorage:  p.TheaterStorage,
		showtimeStorage: p.ShowtimeStorage,
		seatStorage:     p.SeatStorage,
		ticketStorage:   p.TicketStorage,
	}

	return SeatServiceResult{
		SeatService: s,
	}
}

func (s *SeatServiceImpl) GetAvailableSeats(ctx context.Context, req *theater_proto.GetAvailableSeatsRequest) (*theater_proto.GetAvailableSeatsResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	availableSeats, err := s.seatStorage.FindShowtimeAvailableSeats(ctx, &entity.FindShowtimeAvailableSeats{
		ShowtimeID: req.GetShowtimeId(),
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get available seats", zap.Error(err))
		return nil, err
	}

	availableSeatRes := make([]*theater_proto.GetAvailableSeatsResponse_Seat, 0, len(availableSeats))
	for _, seat := range availableSeats {
		availableSeatRes = append(availableSeatRes, &theater_proto.GetAvailableSeatsResponse_Seat{
			SeatId:     seat.SeatID,
			SeatRow:    seat.SeatRow,
			SeatColumn: seat.SeatColumn,
		})
	}

	return &theater_proto.GetAvailableSeatsResponse{
		Seats: availableSeatRes,
	}, nil
}
