package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/harmonify/movie-reservation-system/pkg/cache"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	error_constant "github.com/harmonify/movie-reservation-system/pkg/error/constant"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	OtpRedisRepositoryParam struct {
		fx.In

		Redis  *cache.Redis
		Logger logger.Logger
		Tracer tracer.Tracer
		Util   *util.Util
		Config *config.Config
	}

	OtpRedisRepositoryResult struct {
		fx.Out

		OtpRedisRepository shared.OtpStorage
	}

	otpRedisRepositoryImpl struct {
		redis     *cache.Redis
		logger    logger.Logger
		tracer    tracer.Tracer
		util      *util.Util
		keyPrefix string
	}
)

func NewOtpRedisRepository(p OtpRedisRepositoryParam) OtpRedisRepositoryResult {
	return OtpRedisRepositoryResult{
		OtpRedisRepository: &otpRedisRepositoryImpl{
			redis:     p.Redis,
			logger:    p.Logger,
			tracer:    p.Tracer,
			util:      p.Util,
			keyPrefix: p.Config.ServiceIdentifier,
		},
	}
}

func (r *otpRedisRepositoryImpl) constructCacheKey(ctx context.Context, otpType string, userIdentifier string) (string, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	userIdHash, err := r.util.EncryptionUtil.SHA256Hasher.Hash(userIdentifier)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to hash user identifier", zap.Error(err))
	}
	return fmt.Sprintf("%s:otp:%s:%s", r.keyPrefix, otpType, userIdHash), err
}

func (r *otpRedisRepositoryImpl) SaveEmailVerificationToken(
	ctx context.Context,
	p shared.SaveEmailVerificationTokenParam,
) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "email", p.Email)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct email verification token cache key", zap.Error(err))
		return err
	}

	_, err = r.redis.Client.Set(ctx, cacheKey, p.Token, p.TTL).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to set email verification token", zap.Error(err))
		return err
	}

	return err
}

func (r *otpRedisRepositoryImpl) GetEmailVerificationToken(ctx context.Context, email string) (string, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "email", email)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct email verification token cache key", zap.Error(err))
		return "", err
	}

	token, err := r.redis.Client.Get(ctx, cacheKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", error_constant.ErrNotFound
		}
		r.logger.WithCtx(ctx).Error("Failed to generate email verification token", zap.Error(err))
		return "", err
	}

	return token, nil
}

func (r *otpRedisRepositoryImpl) DeleteEmailVerificationToken(ctx context.Context, email string) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "email", email)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct email verification token cache key", zap.Error(err))
		return false, err
	}

	removed, err := r.redis.Client.Del(ctx, cacheKey).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to delete email verification token", zap.Error(err))
		return false, err
	}

	return removed == 1, nil
}

func (r *otpRedisRepositoryImpl) SavePhoneOtp(ctx context.Context, p shared.SavePhoneOtpParam) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "phone", p.PhoneNumber)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct phone OTP cache key", zap.Error(err))
		return err
	}

	_, err = r.redis.Client.Set(ctx, cacheKey, p.Otp, p.TTL).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to set phone OTP", zap.Error(err))
		return err
	}

	return err
}

func (r *otpRedisRepositoryImpl) GetPhoneOtp(ctx context.Context, phoneNumber string) (string, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "phone", phoneNumber)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct phone OTP cache key", zap.Error(err))
		return "", err
	}

	otp, err := r.redis.Client.Get(ctx, cacheKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", error_constant.ErrNotFound
		}
		r.logger.WithCtx(ctx).Error("Failed to generate phone OTP", zap.Error(err))
		return "", err
	}

	return otp, nil
}

func (r *otpRedisRepositoryImpl) DeletePhoneOtp(ctx context.Context, phoneNumber string) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "phone", phoneNumber)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct phone OTP cache key", zap.Error(err))
		return false, err
	}

	removed, err := r.redis.Client.Del(ctx, cacheKey).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to delete phone OTP", zap.Error(err))
		return false, err
	}

	return removed == 1, nil
}

func (r *otpRedisRepositoryImpl) IncrementPhoneOtpAttempt(ctx context.Context, phoneNumber string) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "phone-attempt", phoneNumber)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct phone OTP attempt cache key", zap.Error(err))
		return err
	}

	_, err = r.redis.Client.Incr(ctx, cacheKey).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to increment phone OTP attempt", zap.Error(err))
		return err
	}

	return nil
}

func (r *otpRedisRepositoryImpl) GetPhoneOtpAttempt(ctx context.Context, phoneNumber string) (int, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "phone-attempt", phoneNumber)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct phone OTP attempt cache key", zap.Error(err))
		return 0, err
	}

	attempt, err := r.redis.Client.Get(ctx, cacheKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, error_constant.ErrNotFound
		}
		r.logger.WithCtx(ctx).Error("Failed to increment phone OTP attempt", zap.Error(err))
		return 0, err
	}

	attemptInt, err := strconv.Atoi(attempt)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to parse phone OTP attempt", zap.Error(err))
		return 0, err
	}

	return attemptInt, nil
}

func (r *otpRedisRepositoryImpl) DeletePhoneOtpAttempt(ctx context.Context, phoneNumber string) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "phone-attempt", phoneNumber)
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct phone OTP attempt cache key", zap.Error(err))
		return false, err
	}

	removed, err := r.redis.Client.Del(ctx, cacheKey).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to delete phone OTP attempt", zap.Error(err))
		return false, err
	}

	return removed == 1, nil
}
