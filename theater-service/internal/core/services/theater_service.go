package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/shared"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	TheaterService interface {
		GetActiveMovies(ctx context.Context, req *theater_proto.GetActiveMoviesRequest) (*theater_proto.GetActiveMoviesResponse, error)
		GetActiveShowtimes(ctx context.Context, req *theater_proto.GetActiveShowtimesRequest) (*theater_proto.GetActiveShowtimesResponse, error)
		GetAvailableSeats(ctx context.Context, req *theater_proto.GetAvailableSeatsRequest) (*theater_proto.GetAvailableSeatsResponse, error)
	}

	TheaterServiceParam struct {
		fx.In
		Logger          logger.Logger
		Tracer          tracer.Tracer
		TheaterStorage  shared.TheaterStorage
		ShowtimeStorage shared.ShowtimeStorage
		SeatStorage     shared.SeatStorage
		TicketStorage   shared.TicketStorage
	}

	TheaterServiceResult struct {
		fx.Out

		TheaterService TheaterService
	}

	theaterServiceImpl struct {
		logger          logger.Logger
		tracer          tracer.Tracer
		theaterStorage  shared.TheaterStorage
		showtimeStorage shared.ShowtimeStorage
		seatStorage     shared.SeatStorage
		ticketStorage   shared.TicketStorage
	}
)

func NewTheaterService(p TheaterServiceParam) TheaterServiceResult {
	s := &theaterServiceImpl{
		logger:          p.Logger,
		tracer:          p.Tracer,
		theaterStorage:  p.TheaterStorage,
		showtimeStorage: p.ShowtimeStorage,
		seatStorage:     p.SeatStorage,
		ticketStorage:   p.TicketStorage,
	}

	return TheaterServiceResult{
		TheaterService: s,
	}
}

func (s *theaterServiceImpl) GetActiveMovies(ctx context.Context, req *theater_proto.GetActiveMoviesRequest) (*theater_proto.GetActiveMoviesResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	theaterId := req.GetTheaterId()
	if theaterId == "" {
		return nil, TheaterIDRequiredError
	}

	findModel := &entity.FindManyShowtimes{
		TheaterID: sql.NullString{String: theaterId, Valid: true},
	}
	if req.GetIncludeUpcoming() {
		findModel.StartTimeGte = sql.NullTime{Time: time.Now(), Valid: true}
		findModel.StartTimeLte = sql.NullTime{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true}
	}

	activeShowtimes, err := s.showtimeStorage.FindManyShowtimes(ctx, findModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get active showtimes", zap.Error(err))
		return nil, err
	}

	activeMovies := make([]*theater_proto.GetActiveMoviesResponse_Movie, 0)
	for _, showtime := range activeShowtimes {
		activeMovies = append(activeMovies, &theater_proto.GetActiveMoviesResponse_Movie{
			MovieId: showtime.MovieID,
		})
	}

	return &theater_proto.GetActiveMoviesResponse{
		Movies: activeMovies,
	}, nil
}

func (s *theaterServiceImpl) GetActiveShowtimes(ctx context.Context, req *theater_proto.GetActiveShowtimesRequest) (*theater_proto.GetActiveShowtimesResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	theaterId := req.GetTheaterId()
	if theaterId == "" {
		return nil, TheaterIDRequiredError
	}

	movieId := req.GetMovieId()
	if movieId == "" {
		return nil, MovieIDRequiredError
	}

	activeShowtimes, err := s.showtimeStorage.FindManyShowtimes(ctx, &entity.FindManyShowtimes{
		TheaterID:    sql.NullString{String: theaterId, Valid: true},
		MovieID:      sql.NullString{String: movieId, Valid: true},
		StartTimeGte: sql.NullTime{Time: time.Now(), Valid: true},
		StartTimeLte: sql.NullTime{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get active showtimes", zap.Error(err))
		return nil, err
	}

	showtimeRoomIds := make([]string, 0, len(activeShowtimes))
	for _, showtime := range activeShowtimes {
		showtimeRoomIds = append(showtimeRoomIds, showtime.RoomID)
	}

	totalSeats, err := s.seatStorage.CountRoomSeats(ctx, showtimeRoomIds)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get total seats", zap.Error(err))
		return nil, err
	}

	totalSeatsMap := make(map[string]uint32, len(totalSeats))
	for _, t := range totalSeats {
		totalSeatsMap[t.RoomID] = t.Count
	}

	activeShowtimeIds := make([]string, 0, len(activeShowtimes))
	for _, showtime := range activeShowtimes {
		activeShowtimeIds = append(activeShowtimeIds, showtime.ShowtimeID)
	}

	showtimeSeats, err := s.ticketStorage.CountShowtimeTickets(ctx, activeShowtimeIds)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get available seats", zap.Error(err))
		return nil, err
	}

	showtimeSeatsMap := make(map[string]uint32, len(showtimeSeats))
	for _, t := range showtimeSeats {
		showtimeSeatsMap[t.ShowtimeID] = t.Count
	}

	showtimes := make([]*theater_proto.GetActiveShowtimesResponse_Showtime, 0, len(activeShowtimes))
	for _, showtime := range activeShowtimes {
		showtimes = append(showtimes, &theater_proto.GetActiveShowtimesResponse_Showtime{
			ShowtimeId:     showtime.ShowtimeID,
			StartTime:      uint32(showtime.StartTime.Unix()),
			AvailableSeats: totalSeatsMap[showtime.RoomID] - showtimeSeatsMap[showtime.ShowtimeID],
		})
	}

	return &theater_proto.GetActiveShowtimesResponse{
		Showtimes: showtimes,
	}, nil
}

func (s *theaterServiceImpl) GetAvailableSeats(ctx context.Context, req *theater_proto.GetAvailableSeatsRequest) (*theater_proto.GetAvailableSeatsResponse, error) {
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
