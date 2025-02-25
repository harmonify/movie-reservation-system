package service

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
	ShowtimeService interface {
		GetActiveMovies(ctx context.Context, req *theater_proto.GetActiveMoviesRequest) (*theater_proto.GetActiveMoviesResponse, error)
		GetActiveShowtimes(ctx context.Context, req *theater_proto.GetActiveShowtimesRequest) (*theater_proto.GetActiveShowtimesResponse, error)
	}

	ShowtimeServiceParam struct {
		fx.In
		Logger          logger.Logger
		Tracer          tracer.Tracer
		TheaterStorage  shared.TheaterStorage
		ShowtimeStorage shared.ShowtimeStorage
		SeatStorage     shared.SeatStorage
		TicketStorage   shared.TicketStorage
	}

	ShowtimeServiceResult struct {
		fx.Out

		ShowtimeService ShowtimeService
	}

	showtimeServiceImpl struct {
		logger          logger.Logger
		tracer          tracer.Tracer
		theaterStorage  shared.TheaterStorage
		showtimeStorage shared.ShowtimeStorage
		seatStorage     shared.SeatStorage
		ticketStorage   shared.TicketStorage
	}
)

func NewShowtimeService(p ShowtimeServiceParam) ShowtimeServiceResult {
	s := &showtimeServiceImpl{
		logger:          p.Logger,
		tracer:          p.Tracer,
		theaterStorage:  p.TheaterStorage,
		showtimeStorage: p.ShowtimeStorage,
		seatStorage:     p.SeatStorage,
		ticketStorage:   p.TicketStorage,
	}

	return ShowtimeServiceResult{
		ShowtimeService: s,
	}
}

func (s *showtimeServiceImpl) GetActiveMovies(ctx context.Context, req *theater_proto.GetActiveMoviesRequest) (*theater_proto.GetActiveMoviesResponse, error) {
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

	res, err := s.showtimeStorage.FindManyShowtimes(ctx, findModel)
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get active showtimes", zap.Error(err))
		return nil, err
	}

	activeMovies := make([]*theater_proto.GetActiveMoviesResponse_Movie, 0)
	for _, showtime := range res.Showtimes {
		activeMovies = append(activeMovies, &theater_proto.GetActiveMoviesResponse_Movie{
			MovieId: showtime.MovieID,
		})
	}

	return &theater_proto.GetActiveMoviesResponse{
		Movies: activeMovies,
	}, nil
}

func (s *showtimeServiceImpl) GetActiveShowtimes(ctx context.Context, req *theater_proto.GetActiveShowtimesRequest) (*theater_proto.GetActiveShowtimesResponse, error) {
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

	res, err := s.showtimeStorage.FindManyShowtimes(ctx, &entity.FindManyShowtimes{
		TheaterID:    sql.NullString{String: theaterId, Valid: true},
		MovieID:      sql.NullString{String: movieId, Valid: true},
		StartTimeGte: sql.NullTime{Time: time.Now(), Valid: true},
		StartTimeLte: sql.NullTime{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("Failed to get active showtimes", zap.Error(err))
		return nil, err
	}

	showtimeRoomIds := make([]string, 0, len(res.Showtimes))
	for _, showtime := range res.Showtimes {
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

	activeShowtimeIds := make([]string, 0, len(res.Showtimes))
	for _, showtime := range res.Showtimes {
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

	showtimes := make([]*theater_proto.GetActiveShowtimesResponse_Showtime, 0, len(res.Showtimes))
	for _, showtime := range res.Showtimes {
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
