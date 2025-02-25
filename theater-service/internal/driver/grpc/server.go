package grpc_driver

import (
	"context"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/service"
	"go.uber.org/fx"
)

func RegisterTheaterServiceServer(
	server *grpc_pkg.GrpcServer,
	handler theater_proto.TheaterServiceServer,
) {
	theater_proto.RegisterTheaterServiceServer(server.Server, handler)
}

type TheaterServiceServerParam struct {
	fx.In
	Logger          logger.Logger
	Tracer          tracer.Tracer
	ErrorMapper     error_pkg.ErrorMapper
	ShowtimeService service.ShowtimeService
	SeatService     service.SeatService
}

type TheaterServiceServerImpl struct {
	theater_proto.UnimplementedTheaterServiceServer // Embedding for compatibility
	logger                                          logger.Logger
	tracer                                          tracer.Tracer
	errorMapper                                     error_pkg.ErrorMapper
	showtimeService                                 service.ShowtimeService
	seatService                                     service.SeatService
}

func NewTheaterServiceServer(
	p TheaterServiceServerParam,
) theater_proto.TheaterServiceServer {
	return &TheaterServiceServerImpl{
		UnimplementedTheaterServiceServer: theater_proto.UnimplementedTheaterServiceServer{},
		logger:                            p.Logger,
		tracer:                            p.Tracer,
		errorMapper:                       p.ErrorMapper,
		showtimeService:                   p.ShowtimeService,
		seatService:                       p.SeatService,
	}
}

func (s *TheaterServiceServerImpl) GetActiveMovies(ctx context.Context, req *theater_proto.GetActiveMoviesRequest) (*theater_proto.GetActiveMoviesResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.showtimeService.GetActiveMovies(ctx, req)
	if err != nil {
		return nil, s.errorMapper.ToGrpcError(err)
	}

	return res, nil
}

func (s *TheaterServiceServerImpl) GetActiveShowtimes(ctx context.Context, req *theater_proto.GetActiveShowtimesRequest) (*theater_proto.GetActiveShowtimesResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.showtimeService.GetActiveShowtimes(ctx, req)
	if err != nil {
		return nil, s.errorMapper.ToGrpcError(err)
	}

	return res, nil
}

func (s *TheaterServiceServerImpl) GetAvailableSeats(ctx context.Context, req *theater_proto.GetAvailableSeatsRequest) (*theater_proto.GetAvailableSeatsResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := s.seatService.GetAvailableSeats(ctx, req)
	if err != nil {
		return nil, s.errorMapper.ToGrpcError(err)
	}

	return res, nil
}
