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
	mem      *map[uint64]uint64
}

func (b *Fibonacci) Duration() time.Duration {
	d := b.ForAttempt(float64(atomic.AddUint64(&b.attempt, 1) - 1))
	return d
}

func (b *Fibonacci) ForAttempt(attempt float64) time.Duration {
	min := b.Min
	if min <= 0 {
		min = defaultMin
	}

	max := b.Max
	if max <= 0 {
		max = defaultMax
	}

	if min >= max {
		return max
	}

	minf := float64(min)
	durf := minf * float64(b.fib(uint64(attempt)))

	if b.Jitter {
		durf = rand.Float64()*(durf-minf) + minf
	}

	if durf > maxInt64 {
		return max
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
	if b.mem == nil {
		b.mem = &map[uint64]uint64{}
	}
	m := *b.mem
	if val, ok := m[n]; ok {
		return val
	}
	if n <= 1 {
		return n
	}
	res := b.fib(n-1) + b.fib(n-2)
	m[n] = res
	return res
}
