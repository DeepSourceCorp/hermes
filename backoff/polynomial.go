package backoff

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

type Polynomial struct {
	attempt  uint64
	Exponent float64
	Jitter   bool
	Min, Max time.Duration
}

const defaultExponent float64 = 1

// NexPolynomial is a constructor for Polynomial.
func NewPolynomial(Exponent float64, Jitter bool, Min, Max time.Duration) *Polynomial {
	if Min <= 0 {
		Min = defaultMin
	}
	if Max <= 0 {
		Max = defaultMax
	}
	if Max > time.Duration(maxInt64) {
		Max = defaultMax
	}
	if Exponent <= 0 {
		Exponent = defaultExponent
	}

	if Min >= Max {
		Min = defaultMin
		Max = defaultMax
	}

	return &Polynomial{
		Exponent: Exponent,
		Min:      Min,
		Max:      Max,
		Jitter:   Jitter,
	}
}

func (b *Polynomial) Duration() time.Duration {
	d := b.ForAttempt(atomic.AddUint64(&b.attempt, 1) - 1)
	return d
}

func (b *Polynomial) ForAttempt(attempt uint64) time.Duration {
	min := b.Min
	max := b.Max
	exponent := b.Exponent

	minf := float64(min)
	attemptf := float64(attempt)
	durf := minf * math.Pow(attemptf, exponent)

	if b.Jitter {
		durf = rand.Float64()*(durf-minf) + minf
	}

	dur := time.Duration(durf)
	if dur > max {
		return max
	}

	return dur
}

func (b *Polynomial) Reset() {
	atomic.StoreUint64(&b.attempt, 0)
}

func (b *Polynomial) Attempt() float64 {
	return float64(atomic.LoadUint64(&b.attempt))
}

func (b *Polynomial) Copy() *Polynomial {
	return &Polynomial{
		Exponent: b.Exponent,
		Jitter:   b.Jitter,
		Min:      b.Min,
		Max:      b.Max,
	}
}
