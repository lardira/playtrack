package tech

import (
	"context"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func TestHealthChecker_CtxDone(t *testing.T) {
	pollInterval := 1 * time.Hour
	pinger := NewMockPinger(t)
	checker := NewHealthChecker(pinger, pollInterval, "test")

	pinger.AssertNotCalled(t, "Ping")

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	checker.Check(ctx)

	assert.True(t, checker.Ok())
}
