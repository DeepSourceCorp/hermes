package backoff

import (
	"testing"
	"time"
)

func TestNewFibonacci(t *testing.T) {
	t.Run("generated duration must generate sanitized object", func(t *testing.T) {
		got := NewFibonacci(true, -1*time.Second, -1*time.Second)

		if got.Min != defaultMin {
			t.Errorf("NewFibonacci() Exponential.Min = %v, want %v", got.Min, defaultMin)
		}

		if got.Max != defaultMax {
			t.Errorf("NewFibonacci() Exponential.Max = %v, want %v", got.Max, defaultMax)
		}

		got = NewFibonacci(false, -1*time.Second, time.Duration(maxInt64)+10)

		if got.Max != defaultMax {
			t.Errorf("NewFibonacci() Exponential.Max = %v, want %v", got.Max, defaultMax)
		}
	})

	t.Run("Min and Max should fallback to default if Min > Max", func(t *testing.T) {
		got := NewFibonacci(true, 20*time.Second, 1*time.Second)
		if got.Max != defaultMax || got.Min != defaultMin {
			t.Errorf(
				"NewFibonacci() Exponential.Max = %v, want %v, Exponential.Min=%v, want = %v",
				got.Max, defaultMax, got.Min, defaultMin,
			)
		}
	})
}

func TestFibonacci_Duration(t *testing.T) {
	type fields struct {
		attempt uint64
		Jitter  bool
		Min     time.Duration
		Max     time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			"must generate the correct duration for attempt = 6",
			fields{
				attempt: uint64(6),
				Min:     1 * time.Second,
				Max:     10 * time.Second,
			},
			8 * time.Second,
		},
		{
			"must return Max if duration value greater than Max",
			fields{
				attempt: uint64(20),
				Min:     1 * time.Second,
				Max:     10 * time.Second,
			},
			10 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Fibonacci{
				attempt: tt.fields.attempt,
				Jitter:  tt.fields.Jitter,
				Min:     tt.fields.Min,
				Max:     tt.fields.Max,
				cache:   map[uint64]uint64{},
			}
			if got := b.Duration(); got != tt.want {
				t.Errorf("Fibonacci.Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFibonacci_Duration_WithJitter(t *testing.T) {
	b := &Fibonacci{
		attempt: uint64(6),
		Jitter:  true,
		Min:     1 * time.Second,
		Max:     10 * time.Second,
		cache:   map[uint64]uint64{},
	}
	// TODO: This is probably not deterministic enough.
	t.Run("must generate a random duration between the minimum and the exponential value", func(t *testing.T) {
		got := b.Duration()
		if got < b.Min || got > 8*time.Second {
			t.Errorf("Exponential.Duration() = %v, want %v > x < %v", got, b.Min, 9*time.Second)
		}
	})
}
