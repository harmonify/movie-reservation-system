package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/harmonify/movie-reservation-system/pkg/cache"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/mennanov/limiters"
	"go.uber.org/fx"
)

const (
	DefaultCapacity   = int64(2)
	DefaultRefillRate = 3 * time.Second
)

type RateLimiterRegistry interface {
	Len() int
	GetHttpRequestRateLimiter(p *HttpRequestRateLimiterParam, c *RateLimiterConfig) (RateLimiter, error)
}

type RateLimiter interface {
	Take(ctx context.Context, tokens int64) (retryAfter time.Duration, err error)
	Limit(ctx context.Context) (retryAfter time.Duration, err error)
	Reset(ctx context.Context) error
}

type RateLimiterRegistryParam struct {
	fx.In

	Logger logger.Logger
	Redis  *cache.Redis
}

type rateLimiterRegistryImpl struct {
	logger        logger.Logger
	wrappedLogger *RateLimiterLogger
	redis         *cache.Redis
	registry      *limiters.Registry
	clock         limiters.Clock
	keyPrefix     string
}

func NewRateLimiterRegistry(p RateLimiterRegistryParam, c *RateLimiterRegistryConfig) (RateLimiterRegistry, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(c); err != nil {
		return nil, err
	}

	registry := limiters.NewRegistry()
	clock := limiters.NewSystemClock()

	go func() {
		// Garbage collect old limiters to prevent memory leaks
		for {
			<-time.After(DefaultRefillRate)
			registry.DeleteExpired(clock.Now())
		}
	}()

	rl := &rateLimiterRegistryImpl{
		logger:        p.Logger,
		wrappedLogger: NewRateLimiterLogger(p.Logger),
		redis:         p.Redis,
		registry:      registry,
		clock:         clock,
		keyPrefix:     c.ServiceIdentifier,
	}

	return rl, nil
}

func (rl *rateLimiterRegistryImpl) Len() int {
	return rl.registry.Len()
}

// GetHttpRequestRateLimiter returns a rate limiter for HTTP requests based on IP address, HTTP method, and HTTP path
func (rl *rateLimiterRegistryImpl) GetHttpRequestRateLimiter(p *HttpRequestRateLimiterParam, c *RateLimiterConfig) (RateLimiter, error) {
	key := rl.constructLimiterKey(p)

	bucket := rl.registry.GetOrCreate(
		key,
		func() interface{} {
			return rl.createLimiter(p, c)
		},
		c.RefillRate,
		rl.clock.Now(),
	)
	tokenBucket, ok := bucket.(*limiters.TokenBucket)
	if !ok {
		return nil, fmt.Errorf("failed to cast bucket to TokenBucket")
	}

	return tokenBucket, nil
}

func (rl *rateLimiterRegistryImpl) constructLimiterKey(p *HttpRequestRateLimiterParam) string {
	return fmt.Sprintf("%s:rl:%s:%s:%s", rl.keyPrefix, p.ID, p.Method, p.Path)
}

func (rl *rateLimiterRegistryImpl) constructLockerKey(p *HttpRequestRateLimiterParam) string {
	return fmt.Sprintf("%s:rll:%s:%s:%s", rl.keyPrefix, p.ID, p.Method, p.Path)
}

func (rl *rateLimiterRegistryImpl) createLimiter(p *HttpRequestRateLimiterParam, c *RateLimiterConfig) *limiters.TokenBucket {
	key := rl.constructLimiterKey(p)
	tokenBucketStateBackend := limiters.NewTokenBucketRedis(
		rl.redis.Client,
		key,
		c.RefillRate,
		false,
	)

	pool := goredis.NewPool(rl.redis.Client)
	lockerKey := rl.constructLockerKey(p)
	locker := limiters.NewLockRedis(pool, lockerKey)

	limiter := limiters.NewTokenBucket(
		c.Capacity,
		c.RefillRate,
		locker,
		tokenBucketStateBackend,
		limiters.NewSystemClock(),
		rl.wrappedLogger,
	)

	return limiter
}
