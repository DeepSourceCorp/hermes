package backoff

import (
	"math/rand"
	"sync/atomic"
	"time"
)

type Linear struct {
	attempt  uint64
	Jitter   bool
	Min, Max time.Duration
}

// NexLinear is a constructor for Linear.
func NewLinear(Jitter bool, Min, Max time.Duration) *Linear {
	if Min <= 0 {
		Min = defaultMin
	}
	if Max <= 0 {
		Max = defaultMax
	}
	if Max > time.Duration(maxInt64) {
		Max = defaultMax
	}

	if Min >= Max {
		Min = defaultMin
		Max = defaultMax
	}

	return &Linear{
		Min:    Min,
		Max:    Max,
		Jitter: Jitter,
	}
}

func (b *Linear) Duration() time.Duration {
	d := b.ForAttempt(atomic.AddUint64(&b.attempt, 1) - 1)
	return d
}

func (b *Linear) ForAttempt(attempt uint64) time.Duration {
	min := b.Min
	max := b.Max

	minf := float64(min)
	attemptf := float64(attempt)
	durf := minf * attemptf

	if b.Jitter {
		durf = rand.Float64()*(durf-minf) + minf
	}

	dur := time.Duration(durf)
	if dur > max {
		return max
	}

	return dur
}

func (b *Linear) Reset() {
	atomic.StoreUint64(&b.attempt, 0)
}

func (b *Linear) Attempt() float64 {
	return float64(atomic.LoadUint64(&b.attempt))
}

func (b *Linear) Copy() *Linear {
	return &Linear{
		Jitter: b.Jitter,
		Min:    b.Min,
		Max:    b.Max,
	}
}
