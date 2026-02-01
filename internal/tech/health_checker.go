package tech

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

const (
	maxConsecutiveErr = 5

	defaultPollInterval = 10 * time.Second
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthChecker struct {
	errCount     atomic.Uint32
	pollInterval time.Duration
	pinger       Pinger
	tag          string
}

func NewHealthChecker(pinger Pinger, pollInterval time.Duration, tag string) *HealthChecker {
	if pollInterval == 0 {
		pollInterval = defaultPollInterval
	}
	return &HealthChecker{
		pinger:       pinger,
		pollInterval: pollInterval,
		tag:          tag,
	}
}

func (c *HealthChecker) Check(ctx context.Context) {
	ticker := time.Tick(c.pollInterval)

	for {
		select {
		case <-ticker:
			if err := c.pinger.Ping(ctx); err != nil {
				log.Printf("could not ping %s", c.tag)

				newCount := c.errCount.Add(1)
				if newCount >= maxConsecutiveErr {
					log.Panicf("could not ping %s for %d times", c.tag, newCount)
				}
				continue
			}

			if c.errCount.Load() > 0 {
				c.errCount.Store(0)
			}

		case <-ctx.Done():
			log.Printf("health checker stopped")
			return
		}
	}
}

func (c *HealthChecker) Ok() bool {
	return c.errCount.Load() < maxConsecutiveErr
}
