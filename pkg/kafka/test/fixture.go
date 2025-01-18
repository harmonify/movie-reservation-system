package test

import "github.com/google/uuid"

var (
	runId           = uuid.New().String()
	TestBasicTopic  = "test_basic_" + runId
	TestRouterTopic = "test_router_" + runId
)
