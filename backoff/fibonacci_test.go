package backoff

import (
	"testing"
	"time"
)

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
			"must generate the correct duration for attempt=6",
			fields{
				attempt: uint64(6),
				Jitter:  false,
				Min:     1 * time.Second,
				Max:     10 * time.Minute,
			},
			8 * time.Second,
		},
		{
			"must return Max if duration greater than Max",
			fields{
				attempt: uint64(30),
				Jitter:  false,
				Min:     1 * time.Minute,
				Max:     2 * time.Minute,
			},
			2 * time.Minute,
		},
		{
			"must use default Min if Min lesser than 0",
			fields{
				attempt: 1,
				Min:     -1 * time.Second,
				Max:     10 * time.Minute,
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
			b := &Fibonacci{
				attempt: tt.fields.attempt,
				Jitter:  tt.fields.Jitter,
				Min:     tt.fields.Min,
				Max:     tt.fields.Max,
			}
			if got := b.Duration(); got != tt.want {
				t.Errorf("Fibonacci.Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFibonacci_fib(t *testing.T) {
	b := Fibonacci{}
	want := uint64(6765)
	if got := b.fib(20); got != want {
		t.Errorf("Fibonacci.fib() = %v, want %v", got, want)
	}

}
