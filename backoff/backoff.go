package backoff

import (
	"math"
	"time"
)

const defaultMin time.Duration = 100 * time.Millisecond
const defaultMax time.Duration = 10 * time.Minute
const maxInt64 = float64(math.MaxInt64 - 512)

type Backoff interface {
	Duration() time.Duration
	ForAttempt(attempt float64) time.Duration
}
