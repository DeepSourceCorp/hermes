package backoff

import (
	"fmt"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

type Exponential struct {
	attempt uint64

	Factor float64

	Jitter bool

	Min, Max time.Duration
}

func (b *Exponential) Duration() time.Duration {
	d := b.ForAttempt(float64(atomic.AddUint64(&b.attempt, 1) - 1))
	return d
}

const maxInt64 = float64(math.MaxInt64 - 512)

func (b *Exponential) ForAttempt(attempt float64) time.Duration {

	min := b.Min
	fmt.Println(min)
	if min <= 0 {
		min = defaultMin
	}

	fmt.Println(min)

	max := b.Max
	if max <= 0 {
		max = defaultMax
	}

	if min >= max {
		return max
	}

	factor := b.Factor
	if factor <= 0 {
		factor = 2
	}

	minf := float64(min)
	durf := minf * math.Pow(factor, attempt)

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
