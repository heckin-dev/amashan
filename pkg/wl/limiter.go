package wl

import (
	"github.com/hashicorp/go-hclog"
	"math"
	"sync"
	"time"
)

type PointLimiter struct {
	l hclog.Logger

	limitPerHour        int
	pointsSpentThisHour float64
	resetInSeconds      int

	timer *time.Timer
	mu    *sync.RWMutex
}

// SetPointsSpent updates the total points spent.
func (p *PointLimiter) SetPointsSpent(data RateLimitData) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.limitPerHour = int(data.LimitPerHour)
	p.pointsSpentThisHour = float64(data.PointsSpentThisHour)

	p.l.Debug("PointLimiter", "limit", p.limitPerHour, "spent", p.pointsSpentThisHour)
}

// CanSpendPoints determines if there are points to spend this hour.
func (p *PointLimiter) CanSpendPoints() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return math.Ceil(p.pointsSpentThisHour) < float64(p.limitPerHour)
}

// tick is the func called by time.AfterFunc when the timer ticks.
func (p *PointLimiter) tick() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pointsSpentThisHour = 0
	p.resetInSeconds = 3600

	p.timer.Stop()
	p.timer.Reset(time.Hour)

	p.l.Debug("PointLimiter ticked")
}

// NewPointLimiter creates a new PointLimiter with the supplied values.
func NewPointLimiter(l hclog.Logger, data RateLimitData) *PointLimiter {
	pl := &PointLimiter{
		l:                   l,
		limitPerHour:        int(data.LimitPerHour),
		pointsSpentThisHour: float64(data.PointsSpentThisHour),
		resetInSeconds:      int(data.PointsResetIn),
		mu:                  &sync.RWMutex{},
	}

	pl.timer = time.AfterFunc(time.Duration(pl.resetInSeconds)*time.Second, pl.tick)

	return pl
}
