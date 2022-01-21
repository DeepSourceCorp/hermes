package backoff

import (
	"testing"
	"time"
)

func TestLinear_Duration(t *testing.T) {
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
			"must generate the correct duration for attempt=2",
			fields{
				attempt: uint64(2),
				Min:     1 * time.Second,
				Max:     10 * time.Second,
			},
			2 * time.Second,
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
			b := &Linear{
				attempt: tt.fields.attempt,
				Jitter:  tt.fields.Jitter,
				Min:     tt.fields.Min,
				Max:     tt.fields.Max,
			}
			if got := b.Duration(); got != tt.want {
				t.Errorf("Linear.Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinear_Duration_WithJitter(t *testing.T) {
	b := &Linear{
		attempt: uint64(2),
		Jitter:  true,
		Min:     1 * time.Second,
		Max:     10 * time.Second,
	}
	// TODO: This is probably not deterministic enough.
	t.Run("must generate a random duration between the minimum and the Linear value", func(t *testing.T) {
		got := b.Duration()
		if got < b.Min || got > 9*time.Second {
			t.Errorf("Linear.Duration() = %v, want %v > x < %v", got, b.Min, 9*time.Second)
		}
	})
}

func TestNewLinear(t *testing.T) {
	t.Run("generated duration must generate sanitized object", func(t *testing.T) {
		got := NewLinear(true, -1*time.Second, -1*time.Second)
		if got.Min != defaultMin {
			t.Errorf("NewLinear() Linear.Min = %v, want %v", got.Min, defaultMin)
		}

		if got.Max != defaultMax {
			t.Errorf("NewLinear() Linear.Max = %v, want %v", got.Max, defaultMax)
		}

		got = NewLinear(false, -1*time.Second, time.Duration(maxInt64)+10)

		if got.Max != defaultMax {
			t.Errorf("NewLinear() Linear.Max = %v, want %v", got.Max, defaultMax)
		}
	})

	t.Run("Min and Max should fallback to default if Min > Max", func(t *testing.T) {
		got := NewLinear(true, 20*time.Second, 1*time.Second)
		if got.Max != defaultMax || got.Min != defaultMin {
			t.Errorf(
				"NewLinear() Linear.Max = %v, want %v, Linear.Min=%v, want = %v",
				got.Max, defaultMax, got.Min, defaultMin,
			)
		}
	})
}
