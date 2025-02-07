package ratelimiter

import "time"

type (
	// Default is 2 tokens per 3 seconds
	RateLimiterRegistryConfig struct {
		ServiceIdentifier string `validate:"required"`
	}

	RateLimiterConfig struct {
		// The maximum number of tokens that can be stored in the bucket
		Capacity int64 `validate:"required,min=1"`
		// The rate is used to define the token bucket refill rate (1 token per duration)
		// and also the TTL for the limiters (both in Redis and in the registry).
		RefillRate time.Duration `validate:"required,min=1s"`
	}

	HttpRequestRateLimiterParam struct {
		ID     string // IP address or any unique user identifier
		Method string
		Path   string
	}
)
