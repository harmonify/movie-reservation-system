package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/harmonify/movie-reservation-system/pkg/cache"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	otp_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/otp"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
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
		Config *config.UserServiceConfig
	}

	OtpRedisRepositoryResult struct {
		fx.Out

		OtpRedisRepository shared.OtpCache
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

func (r *otpRedisRepositoryImpl) SaveEmailVerificationCode(
	ctx context.Context,
	p shared.SaveEmailVerificationCodeParam,
) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "email", p.Email)
	r.logger.WithCtx(ctx).Debug("Cache key to save", zap.String("cacheKey", cacheKey))
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct email verification token cache key", zap.Error(err))
		return err
	}

	_, err = r.redis.Client.Set(ctx, cacheKey, p.Code, p.TTL).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to set email verification token", zap.Error(err))
		return err
	}

	return err
}

func (r *otpRedisRepositoryImpl) GetEmailVerificationCode(ctx context.Context, email string) (string, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "email", email)
	r.logger.WithCtx(ctx).Debug("Cache key to get", zap.String("cacheKey", cacheKey))
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct email verification token cache key", zap.Error(err))
		return "", err
	}

	token, err := r.redis.Client.Get(ctx, cacheKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", otp_service.VerificationTokenNotFoundError
		}
		r.logger.WithCtx(ctx).Error("Failed to generate email verification token", zap.Error(err))
		return "", err
	}

	return token, nil
}

func (r *otpRedisRepositoryImpl) DeleteEmailVerificationCode(ctx context.Context, email string) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cacheKey, err := r.constructCacheKey(ctx, "email", email)
	r.logger.WithCtx(ctx).Debug("Cache key to delete", zap.String("cacheKey", cacheKey))
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to construct email verification token cache key", zap.Error(err))
		return false, err
	}

	removed, err := r.redis.Client.Del(ctx, cacheKey).Result()
	if err != nil {
		r.logger.WithCtx(ctx).Error("Failed to delete email verification token", zap.Error(err))
		return false, err
	}

	return removed >= 1, nil
}

func (r *otpRedisRepositoryImpl) SavePhoneNumberVerificationOtp(ctx context.Context, p shared.SavePhoneNumberVerificationOtpParam) error {
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

func (r *otpRedisRepositoryImpl) GetPhoneNumberVerificationOtp(ctx context.Context, phoneNumber string) (string, error) {
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
			return "", otp_service.OtpNotFoundError
		}
		r.logger.WithCtx(ctx).Error("Failed to generate phone OTP", zap.Error(err))
		return "", err
	}

	return otp, nil
}

func (r *otpRedisRepositoryImpl) DeletePhoneNumberVerificationOtp(ctx context.Context, phoneNumber string) (bool, error) {
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

	return removed >= 1, nil
}

func (r *otpRedisRepositoryImpl) IncrementPhoneNumberVerificationAttempt(ctx context.Context, phoneNumber string) error {
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

func (r *otpRedisRepositoryImpl) GetPhoneNumberVerificationAttempt(ctx context.Context, phoneNumber string) (int, error) {
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
			return 0, otp_service.OtpNotFoundError
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

func (r *otpRedisRepositoryImpl) DeletePhoneNumberVerificationAttempt(ctx context.Context, phoneNumber string) (bool, error) {
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

	return removed >= 1, nil
}
