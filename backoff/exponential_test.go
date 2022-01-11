package backoff

import (
	"testing"
	"time"
)

func TestExponential_Duration(t *testing.T) {
	type fields struct {
		attempt uint64
		Factor  float64
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
			"must generate the correct duration for attempt=3",
			fields{
				attempt: uint64(2),
				Factor:  float64(3),
				Jitter:  false,
				Min:     1 * time.Second,
				Max:     10 * time.Second,
			},
			9 * time.Second,
		},
		{
			"must use factor 2 if factor < 0",
			fields{
				attempt: uint64(2),
				Factor:  float64(-1),
				Jitter:  false,
				Min:     1 * time.Second,
				Max:     10 * time.Second,
			},
			4 * time.Second,
		},
		{
			"must return Max if duration greater than Max",
			fields{
				attempt: uint64(3),
				Jitter:  false,
				Min:     1 * time.Minute,
				Max:     2 * time.Minute,
				Factor:  2,
			},
			2 * time.Minute,
		},
		{
			"must use default Min if Min lesser than 0",
			fields{
				attempt: 1,
				Min:     -1 * time.Second,
				Max:     10 * time.Minute,
				Factor:  1,
			},
			defaultMin,
		},
		{
			"must use default Max if Max lesser than 0",
			fields{
				attempt: 1,
				Min:     10 * time.Minute,
				Max:     -1 * time.Second,
			},
			defaultMax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Exponential{
				attempt: tt.fields.attempt,
				Factor:  tt.fields.Factor,
				Jitter:  tt.fields.Jitter,
				Min:     tt.fields.Min,
				Max:     tt.fields.Max,
			}
			if got := b.Duration(); got != tt.want {
				t.Errorf("Exponential.Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
