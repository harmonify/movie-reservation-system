package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/harmonify/movie-reservation-system/pkg/cache"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	otp_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/otp"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var otpCodeField = "code"
var otpAttemptsField = "attempts"

type otpV2RedisRepositoryImpl struct {
	redis     *cache.Redis
	logger    logger.Logger
	tracer    tracer.Tracer
	util      *util.Util
	keyPrefix string
}

func NewOtpV2RedisRepository(p OtpRedisRepositoryParam) shared.OtpCacheV2 {
	return &otpV2RedisRepositoryImpl{
		redis:     p.Redis,
		logger:    p.Logger,
		tracer:    p.Tracer,
		util:      p.Util,
		keyPrefix: p.Config.ServiceIdentifier,
	}
}

func (r *otpV2RedisRepositoryImpl) constructOtpCacheKey(ctx context.Context, uuid string, otpType string) (string, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if otpType == "" {
		return "", errors.New("otpType is required")
	}

	if uuid == "" {
		return "", errors.New("uuid is required")
	}

	uuidHash, err := r.util.EncryptionUtil.SHA256Hasher.Hash(uuid)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to hash user identifier", zap.Error(err))
		return "", error_pkg.InternalServerError
	}

	return fmt.Sprintf("%s:otp:%s:%s", r.keyPrefix, uuidHash, otpType), nil
}

func (r *otpV2RedisRepositoryImpl) SaveOtp(ctx context.Context, uuid string, otpType shared.OtpType, code string) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructOtpCacheKey(ctx, uuid, otpType.Name)
	if err != nil {
		return err
	}

	p := r.redis.Client.Pipeline()
	p.HSet(ctx, cacheKey, otpCodeField, code, otpAttemptsField, 0)
	p.Expire(ctx, cacheKey, otpType.TTL)
	_, err = p.Exec(ctx)
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to save otp", zap.Error(err), zap.String("uuid", uuid), zap.String("otpType", otpType.Name), zap.String(otpCodeField, code), zap.Duration("ttl", otpType.TTL))
		return error_pkg.InternalServerError
	}

	return nil
}

func (r *otpV2RedisRepositoryImpl) GetOtp(ctx context.Context, uuid string, otpType shared.OtpType) (*shared.Otp, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructOtpCacheKey(ctx, uuid, otpType.Name)
	if err != nil {
		return nil, err
	}

	res, err := r.redis.Client.HMGet(ctx, cacheKey, otpCodeField, otpAttemptsField).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, otp_service.OtpNotFoundError
		}
		r.logger.WithCtx(ctx).Error("failed to get otp", zap.Error(err), zap.String("uuid", uuid), zap.String("otpType", otpType.Name))
		return nil, error_pkg.InternalServerError
	}

	if len(res) == 0 || res[0] == nil {
		return nil, otp_service.OtpNotFoundError
	}

	code, ok := res[0].(string)
	if !ok {
		r.logger.WithCtx(ctx).Error("failed to parse otp code", zap.Any(otpCodeField, res[0]))
		return nil, error_pkg.InternalServerError
	}

	attempts, ok := res[1].(string)
	if !ok {
		r.logger.WithCtx(ctx).Error("failed to parse otp attempts", zap.Any(otpAttemptsField, res[1]))
		return nil, error_pkg.InternalServerError
	}

	attemptsInt, err := strconv.Atoi(attempts)
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to convert otp attempts to int", zap.Error(err))
		return nil, error_pkg.InternalServerError
	}

	return &shared.Otp{
		Code:     code,
		Attempts: attemptsInt,
	}, nil
}

func (r *otpV2RedisRepositoryImpl) DeleteOtp(ctx context.Context, uuid string, otpType shared.OtpType) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructOtpCacheKey(ctx, uuid, otpType.Name)
	if err != nil {
		return false, err
	}

	deleted, err := r.redis.Client.Del(ctx, cacheKey).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to delete otp", zap.Error(err))
		return false, error_pkg.InternalServerError
	}

	if deleted > 1 {
		r.logger.WithCtx(ctx).Warn("deleted more than one otp", zap.Int64("deleted", deleted), zap.String("uuid", uuid), zap.String("otpType", otpType.Name))
	}

	return deleted >= 1, nil
}

func (r *otpV2RedisRepositoryImpl) IncrementOtpAttempt(ctx context.Context, uuid string, otpType shared.OtpType) (*shared.Otp, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructOtpCacheKey(ctx, uuid, otpType.Name)
	if err != nil {
		return nil, err
	}

	otp, err := r.GetOtp(ctx, uuid, otpType)
	if err != nil {
		return nil, err
	}

	res, err := r.redis.Client.HIncrBy(ctx, cacheKey, otpAttemptsField, 1).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("failed to increment otp attempts", zap.Error(err), zap.String("uuid", uuid), zap.String("otpType", otpType.Name))
		return nil, error_pkg.InternalServerError
	}

	return &shared.Otp{
		Code:     otp.Code,
		Attempts: int(res),
	}, nil
}
