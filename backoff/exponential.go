package backoff

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

type Exponential struct {
	attempt  uint64
	Factor   float64
	Jitter   bool
	Min, Max time.Duration
}

const defaultFactor float64 = 2

// NexExponential is a constructor for Exponential.
func NewExponential(Factor float64, Jitter bool, Min, Max time.Duration) *Exponential {
	if Min <= 0 {
		Min = defaultMin
	}
	if Max <= 0 {
		Max = defaultMax
	}
	if Max > time.Duration(maxInt64) {
		Max = defaultMax
	}
	if Factor <= 0 {
		Factor = defaultFactor
	}

	if Min >= Max {
		Min = defaultMin
		Max = defaultMax
	}

	return &Exponential{
		Factor: Factor,
		Min:    Min,
		Max:    Max,
		Jitter: Jitter,
	}
}

func (b *Exponential) Duration() time.Duration {
	d := b.ForAttempt(atomic.AddUint64(&b.attempt, 1) - 1)
	return d
}

func (b *Exponential) ForAttempt(attempt uint64) time.Duration {
	min := b.Min
	max := b.Max
	factor := b.Factor

	minf := float64(min)
	attemptf := float64(attempt)
	durf := minf * math.Pow(factor, attemptf)

	if b.Jitter {
		durf = rand.Float64()*(durf-minf) + minf
	}

	dur := time.Duration(durf)
	if dur > max {
		return max
	}

	return dur
}

func (b *Exponential) Reset() {
	atomic.StoreUint64(&b.attempt, 0)
}

func (b *Exponential) Attempt() float64 {
	return float64(atomic.LoadUint64(&b.attempt))
}

func (b *Exponential) Copy() *Exponential {
	return &Exponential{
		Factor: b.Factor,
		Jitter: b.Jitter,
		Min:    b.Min,
		Max:    b.Max,
	}
}
