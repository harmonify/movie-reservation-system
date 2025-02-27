package grpc_driver

import (
	"context"

	movie_service "github.com/harmonify/movie-reservation-system/movie-service/internal/core/service/movie"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	movie_proto "github.com/harmonify/movie-reservation-system/pkg/proto/movie"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
)

func RegisterMovieServiceServer(
	server *grpc_pkg.GrpcServer,
	handler movie_proto.MovieServiceServer,
) {
	movie_proto.RegisterMovieServiceServer(server.Server, handler)
}

type MovieServiceServerParam struct {
	fx.In
	Logger       logger.Logger
	Tracer       tracer.Tracer
	ErrorMapper  error_pkg.ErrorMapper
	MovieService movie_service.MovieService
}

type MovieServiceServerImpl struct {
	movie_proto.UnimplementedMovieServiceServer // Embedding for compatibility
	logger                                      logger.Logger
	tracer                                      tracer.Tracer
	errorMapper                                 error_pkg.ErrorMapper
	movieService                                movie_service.MovieService
}

func NewMovieServiceServer(
	p MovieServiceServerParam,
) movie_proto.MovieServiceServer {
	return &MovieServiceServerImpl{
		UnimplementedMovieServiceServer: movie_proto.UnimplementedMovieServiceServer{},
		logger:                          p.Logger,
		tracer:                          p.Tracer,
		errorMapper:                     p.ErrorMapper,
		movieService:                    p.MovieService,
	}
}

func (s *MovieServiceServerImpl) GetMovieByID(
	ctx context.Context,
	req *movie_proto.GetMovieByIDRequest,
) (*movie_proto.GetMovieByIDResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	movie, err := s.movieService.GetMovieByID(ctx, req.GetMovieId())
	if err != nil {
		return nil, s.errorMapper.ToGrpcError(err)
	}

	return &movie_proto.GetMovieByIDResponse{
		Movie: &movie_proto.Movie{
			MovieId: movie.MovieID,
			Title:   movie.Title,
		},
	}, nil
}
