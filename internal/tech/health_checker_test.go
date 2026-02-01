package tech

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/stretchr/testify/mock"
)

func TestHealthChecker(t *testing.T) {
	pollInterval := 1 * time.Millisecond

	pinger := NewMockPinger(t)
	checker := NewHealthChecker(pinger, pollInterval, "test")

	pinger.On("Ping", mock.Anything).
		Times(3).Return(nil)
	pinger.On("Ping", mock.Anything).
		Times(maxConsecutiveErr - 1).Return(errors.New(""))
	pinger.On("Ping", mock.Anything).
		Times(2).Return(nil)
	pinger.On("Ping", mock.Anything).
		Times(maxConsecutiveErr).Return(errors.New(""))

	assert.Panics(t, func() {
		checker.Check(t.Context())
	})
}

func TestHealthChecker_CtxDone(t *testing.T) {
	pollInterval := 1 * time.Hour
	pinger := NewMockPinger(t)
	checker := NewHealthChecker(pinger, pollInterval, "test")

	pinger.AssertNotCalled(t, "Ping")

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	assert.NotPanics(t, func() {
		checker.Check(ctx)
	})
}

func TestHealthChecker_Ok(t *testing.T) {
	pollInterval := 1 * time.Millisecond
	pingsBeforeCheck := maxConsecutiveErr - 1

	pinger := NewMockPinger(t)
	checker := NewHealthChecker(pinger, pollInterval, "test")

	checkChan := make(chan time.Time)

	pinger.On("Ping", mock.Anything).
		Times(pingsBeforeCheck).Return(errors.New(""))
	pinger.On("Ping", mock.Anything).
		Times(1).WaitUntil(checkChan).Return(errors.New(""))

	go func() {
		defer close(checkChan)

		// waiting for the last return
		time.Sleep(time.Duration(pingsBeforeCheck)*pollInterval + 1)
		assert.True(t, checker.Ok())
	}()

	assert.Panics(t, func() {
		checker.Check(t.Context())
	})
}
