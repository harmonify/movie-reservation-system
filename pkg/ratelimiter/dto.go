package ratelimiter

import "time"

type (
	// Default is 2 tokens per 3 seconds
	RateLimiterConfig struct {
		ServiceIdentifier string `validate:"required"`
		// The maximum number of tokens that can be stored in the bucket
		Capacity int64 `validate:"required,min=1"`
		// The rate is used to define the token bucket refill rate (1 token per duration)
		// and also the TTL for the limiters (both in Redis and in the registry).
		RefillRate time.Duration `validate:"required,min=1s"`
	}
	HttpRequestRateLimiterParam struct {
		IP     string
		Method string
		Path   string
	}
)
