package wl

import (
	"errors"
	"github.com/hashicorp/go-hclog"
	"math"
	"net/http"
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

func (p *PointLimiter) SpendAllPoints() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pointsSpentThisHour = float64(p.limitPerHour)
}

// SetPointsSpent updates the total points spent.
func (p *PointLimiter) SetPointsSpent(data RateLimitData) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.limitPerHour = int(data.LimitPerHour)
	p.pointsSpentThisHour = float64(data.PointsSpentThisHour)

	p.l.Info("PointLimiter", "limit", p.limitPerHour, "spent", p.pointsSpentThisHour)
}

// CanSpendPoints determines if there are points to spend this hour.
func (p *PointLimiter) CanSpendPoints() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if math.Ceil(p.pointsSpentThisHour) < float64(p.limitPerHour) {
		return nil
	}

	return &ErrNoPointsLeft{
		StatusCode:       http.StatusServiceUnavailable,
		RemainingSeconds: p.resetInSeconds,
		Err:              errors.New("no points available to spend"),
	}
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
