package wl

import (
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPointLimiter(t *testing.T) {
	l := hclog.Default()
	l.SetLevel(hclog.Debug)

	pl := NewPointLimiter(l, RateLimitData{
		LimitPerHour:        5,
		PointsSpentThisHour: 0,
		PointsResetIn:       3,
	})

	assert.Nil(t, pl.CanSpendPoints())

	pl.SetPointsSpent(RateLimitData{
		LimitPerHour:        5,
		PointsSpentThisHour: 5,
		PointsResetIn:       3,
	})

	assert.NotNil(t, pl.CanSpendPoints())

	//time.Sleep(5 * time.Second)
	pl.tick()

	assert.Nil(t, pl.CanSpendPoints())
	assert.Equal(t, float64(0), pl.pointsSpentThisHour)
	assert.Equal(t, 3600, pl.resetInSeconds)
}
