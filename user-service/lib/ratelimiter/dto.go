package ratelimiter

import "time"

type (
	// Default is 2 tokens per 3 seconds
	RateLimiterConfig struct {
		// The maximum number of tokens that can be stored in the bucket
		Capacity int64
		// The rate is used to define the token bucket refill rate (1 token per duration)
		// and also the TTL for the limiters (both in Redis and in the registry).
		RefillRate time.Duration
	}
	HttpRequestRateLimiterParam struct {
		IP     string
		Method string
		Path   string
	}
)
