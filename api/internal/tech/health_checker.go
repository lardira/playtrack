package tech

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

const (
	maxHealthyConsecutive = 5
	maxConsecutive        = maxHealthyConsecutive * 5

	defaultPollInterval = 10 * time.Second
	defaultBackoff      = 2 * time.Minute
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type HealthChecker struct {
	errCount     atomic.Uint32
	pollInterval time.Duration
	pinger       Pinger
	tag          string

	onDone func()
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
			currErrCount := c.errCount.Load()
			if currErrCount >= maxConsecutive {
				select {
				case <-time.After(defaultBackoff):
				case <-ctx.Done():
					log.Printf("health checker stopped")
					c.errCount.Store(0)
					return
				}
			}

			err := c.pinger.Ping(ctx)
			if err != nil {
				log.Printf("could not ping %s", c.tag)
				c.errCount.Add(1)
				continue
			}

			if c.errCount.Load() > 0 {
				c.errCount.Store(0)
			}

		case <-ctx.Done():
			log.Printf("health checker stopped")
			c.errCount.Store(0)
			return
		}
	}
}

func (c *HealthChecker) Ok() bool {
	return c.errCount.Load() < maxHealthyConsecutive
}
