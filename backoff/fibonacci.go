package backoff

import (
	"math/rand"
	"sync/atomic"
	"time"
)

type Fibonacci struct {
	attempt  uint64
	Jitter   bool
	Min, Max time.Duration
	cache    map[uint64]uint64 // cache used for memoization
}

func NewFibonacci(Jitter bool, Min, Max time.Duration) *Fibonacci {
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

	return &Fibonacci{
		Min:    Min,
		Max:    Max,
		Jitter: Jitter,
		cache:  map[uint64]uint64{},
	}
}

func (b *Fibonacci) Duration() time.Duration {
	d := b.ForAttempt(atomic.AddUint64(&b.attempt, 1) - 1)
	return d
}

func (b *Fibonacci) ForAttempt(attempt uint64) time.Duration {
	min := b.Min
	max := b.Max

	minf := float64(min)
	durf := minf * float64(b.fib(uint64(attempt)))

	if b.Jitter {
		durf = rand.Float64()*(durf-minf) + minf
	}
	dur := time.Duration(durf)

	if dur > max {
		return max
	}
	return dur
}

func (b *Fibonacci) Reset() {
	atomic.StoreUint64(&b.attempt, 0)
}

func (b *Fibonacci) Attempt() float64 {
	return float64(atomic.LoadUint64(&b.attempt))
}

func (b *Fibonacci) Copy() *Fibonacci {
	return &Fibonacci{
		Jitter: b.Jitter,
		Min:    b.Min,
		Max:    b.Max,
	}
}

func (b Fibonacci) fib(n uint64) uint64 {
	if val, ok := b.cache[n]; ok {
		return val
	}
	if n <= 1 {
		return n
	}
	res := b.fib(n-1) + b.fib(n-2)
	b.cache[n] = res
	return res
}
