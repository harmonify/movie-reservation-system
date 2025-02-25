package grpc_driver

import (
	"context"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	auth_proto "github.com/harmonify/movie-reservation-system/pkg/proto/auth"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	user_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/user"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/rpc/opa"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RegisterAuthServiceServer(
	server *grpc_pkg.GrpcServer,
	handler auth_proto.AuthServiceServer,
) {
	auth_proto.RegisterAuthServiceServer(server.Server, handler)
}

type AuthServiceServerParam struct {
	fx.In
	logger.Logger
	tracer.Tracer
	error_pkg.ErrorMapper
	jwt_util.JwtUtil
	auth_service.AuthService
	user_service.UserService
	opa.OpaClient
}

type AuthServiceServerImpl struct {
	auth_proto.UnimplementedAuthServiceServer // Embedding for compatibility
	logger                                    logger.Logger
	tracer                                    tracer.Tracer
	errorMapper                               error_pkg.ErrorMapper
	jwtUtil                                   jwt_util.JwtUtil
	authService                               auth_service.AuthService
	userService                               user_service.UserService
	opaClient                                 opa.OpaClient
}

func NewAuthServiceServer(
	p AuthServiceServerParam,
) auth_proto.AuthServiceServer {
	return &AuthServiceServerImpl{
		UnimplementedAuthServiceServer: auth_proto.UnimplementedAuthServiceServer{},
		logger:                         p.Logger,
		tracer:                         p.Tracer,
		errorMapper:                    p.ErrorMapper,
		jwtUtil:                        p.JwtUtil,
		authService:                    p.AuthService,
		userService:                    p.UserService,
		opaClient:                      p.OpaClient,
	}
}

func (s *AuthServiceServerImpl) Auth(ctx context.Context, req *auth_proto.AuthRequest) (*auth_proto.AuthResponse, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	payload, err := s.jwtUtil.JWTVerify(ctx, req.GetAccessToken())
	if err != nil {
		s.logger.WithCtx(ctx).Debug("failed to verify jwt", zap.Error(err))
		return &auth_proto.AuthResponse{}, s.errorMapper.ToGrpcError(err)
	}

	// TODO: cache user info
	userInfo, err := s.userService.GetUser(ctx, user_service.GetUserParam{
		UUID: payload.UUID,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("failed to get user", zap.Error(err))
		return &auth_proto.AuthResponse{}, s.errorMapper.ToGrpcError(err)
	}

	var deletedAt *timestamppb.Timestamp
	if userInfo.DeletedAt != nil {
		deletedAt = timestamppb.New(*userInfo.DeletedAt)
	}

	authResult := &auth_proto.AuthResponse{
		UserInfo: &auth_proto.UserInfo{
			Uuid:                  userInfo.UUID,
			Username:              userInfo.Username,
			Email:                 userInfo.Email,
			PhoneNumber:           userInfo.PhoneNumber,
			FirstName:             userInfo.FirstName,
			LastName:              userInfo.LastName,
			IsEmailVerified:       userInfo.IsEmailVerified,
			IsPhoneNumberVerified: userInfo.IsPhoneNumberVerified,
			CreatedAt:             timestamppb.New(userInfo.CreatedAt),
			UpdatedAt:             timestamppb.New(userInfo.UpdatedAt),
			DeletedAt:             deletedAt,
			Roles:                 userInfo.Roles,
		},
	}

	if req.GetPolicyId() == "" {
		authorized, err := s.opaClient.HasAccess(ctx, req.GetPolicyId(), opa.OpaRequestBody{
			Input: userInfo,
		})
		if err != nil {
			s.logger.WithCtx(ctx).Error("failed to check policy", zap.Error(err))
			return authResult, s.errorMapper.ToGrpcError(err)
		}
		if !authorized {
			s.logger.WithCtx(ctx).Debug("authorization failed", zap.Any("authResult", authResult))
			return authResult, s.errorMapper.ToGrpcError(error_pkg.ForbiddenError)
		}
	}

	s.logger.WithCtx(ctx).Debug("authentication succeed", zap.Any("authResult", authResult))

	return authResult, nil
}
