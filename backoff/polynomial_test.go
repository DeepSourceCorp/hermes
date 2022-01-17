package backoff

import (
	"testing"
	"time"
)

func TestPolynomial_Duration(t *testing.T) {
	type fields struct {
		attempt  uint64
		Exponent float64
		Jitter   bool
		Min      time.Duration
		Max      time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			"must generate the correct duration for attempt=3",
			fields{
				attempt:  uint64(2),
				Exponent: float64(3),
				Min:      1 * time.Second,
				Max:      10 * time.Second,
			},
			8 * time.Second,
		},
		{
			"must return Max if duration value greater than Max",
			fields{
				attempt:  uint64(20),
				Exponent: float64(2),
				Min:      1 * time.Second,
				Max:      10 * time.Second,
			},
			10 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Polynomial{
				attempt:  tt.fields.attempt,
				Exponent: tt.fields.Exponent,
				Jitter:   tt.fields.Jitter,
				Min:      tt.fields.Min,
				Max:      tt.fields.Max,
			}
			if got := b.Duration(); got != tt.want {
				t.Errorf("Polynomial.Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPolynomial_Duration_WithJitter(t *testing.T) {
	b := &Polynomial{
		attempt:  uint64(2),
		Exponent: float64(3),
		Jitter:   true,
		Min:      1 * time.Second,
		Max:      10 * time.Second,
	}
	// TODO: This is probably not deterministic enough.
	t.Run("must generate a random duration between the minimum and the Polynomial value", func(t *testing.T) {
		got := b.Duration()
		if got < b.Min || got > 9*time.Second {
			t.Errorf("Polynomial.Duration() = %v, want %v > x < %v", got, b.Min, 9*time.Second)
		}
	})
}

func TestNewPolynomial(t *testing.T) {
	t.Run("generated duration must generate sanitized object", func(t *testing.T) {
		got := NewPolynomial(0, true, -1*time.Second, -1*time.Second)
		if got.Exponent != 1 {
			t.Errorf("NewPolynomial() Polynomial.Exponent = %v, want %v", got.Exponent, defaultExponent)
		}

		if got.Min != defaultMin {
			t.Errorf("NewPolynomial() Polynomial.Min = %v, want %v", got.Min, defaultMin)
		}

		if got.Max != defaultMax {
			t.Errorf("NewPolynomial() Polynomial.Max = %v, want %v", got.Max, defaultMax)
		}

		got = NewPolynomial(2, false, -1*time.Second, time.Duration(maxInt64)+10)

		if got.Max != defaultMax {
			t.Errorf("NewPolynomial() Polynomial.Max = %v, want %v", got.Max, defaultMax)
		}
	})

	t.Run("Min and Max should fallback to default if Min > Max", func(t *testing.T) {
		got := NewPolynomial(0, true, 20*time.Second, 1*time.Second)
		if got.Max != defaultMax || got.Min != defaultMin {
			t.Errorf(
				"NewPolynomial() Polynomial.Max = %v, want %v, Polynomial.Min=%v, want = %v",
				got.Max, defaultMax, got.Min, defaultMin,
			)
		}
	})
}
