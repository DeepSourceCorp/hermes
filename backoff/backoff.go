package backoff

import "time"

const defaultMin time.Duration = 100 * time.Millisecond
const defaultMax time.Duration = 10 * time.Minute
