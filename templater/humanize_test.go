package templater

import (
	"testing"
)

func TestConcatenateWords(t *testing.T) {
	type args struct {
		words       []interface{}
		conjunction string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "concatenates empty array",
			args: args{words: []interface{}{}, conjunction: "and"},
			want: "",
		},
		{
			name: "concatenates array with single element",
			args: args{words: []interface{}{"first"}, conjunction: "and"},
			want: "first",
		},
		{
			name: "concatenates array with multiple elements",
			args: args{words: []interface{}{"first", "second", "third"}, conjunction: "and"},
			want: "first, second and third",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConcatenateWords(tt.args.words, tt.args.conjunction); got != tt.want {
				t.Errorf("ConcatenateWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	type args struct {
		seconds float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "s second",
			args: args{seconds: 1.0},
			want: "1 second",
		},
		{
			name: "less than a minute",
			args: args{seconds: 33.0},
			want: "33 seconds",
		},
		{
			name: "a minute",
			args: args{seconds: 60.0},
			want: "1 minute",
		},
		{
			name: "less than a hour",
			args: args{seconds: 2400.0},
			want: "40 minutes",
		},
		{
			name: "a few seconds & minutes",
			args: args{seconds: 2450.0},
			want: "40 minutes 50 seconds",
		},
		{
			name: "an hour",
			args: args{seconds: 3600.0},
			want: "1 hour",
		},
		{
			name: "less than a day",
			args: args{seconds: 21600.0},
			want: "6 hours",
		},
		{
			name: "a few seconds, minutes & hours",
			args: args{seconds: 4850.0},
			want: "1 hours 20 minutes 50 seconds",
		},
		{
			name: "a day",
			args: args{seconds: 86400.0},
			want: "1 day",
		},
		{
			name: "a couple of days",
			args: args{seconds: 259200.0},
			want: "3 days",
		},
		{
			name: "a few seconds, minutes & hours",
			args: args{seconds: 268200.0},
			want: "3 days 2 hours 30 minutes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Duration(tt.args.seconds); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestPlural(t *testing.T) {
	type args struct {
		quantity float64
		singular string
		plural   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns singular word when value is singular",
			args: args{quantity: 1.0, singular: "singular", plural: "plural"},
			want: "1 singular",
		},
		{
			name: "returns plural word when value is singular",
			args: args{quantity: 2.0, singular: "singular", plural: "plural"},
			want: "2 plural",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Plural(tt.args.quantity, tt.args.singular, tt.args.plural); got != tt.want {
				t.Errorf("Plural() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestPluralWord(t *testing.T) {
	type args struct {
		quantity float64
		singular string
		plural   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns singular word when value is singular",
			args: args{quantity: 1.0, singular: "singular", plural: "plural"},
			want: "singular",
		},
		{
			name: "returns plural word when value is singular",
			args: args{quantity: 2.0, singular: "singular", plural: "plural"},
			want: "plural",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PluralWord(tt.args.quantity, tt.args.singular, tt.args.plural); got != tt.want {
				t.Errorf("PluralWord() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestTruncateQuantity(t *testing.T) {
	type args struct {
		quantity float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "doesn't truncate quantities less than 1000",
			args: args{quantity: 545.0},
			want: "545",
		},
		{
			name: "truncates quantity 1000",
			args: args{quantity: 1000.0},
			want: "1K",
		},
		{
			name: "truncates quantity 1100",
			args: args{quantity: 1100.0},
			want: "1.1K",
		},
		{
			name: "truncates quantities more than 1100",
			args: args{quantity: 34452.0},
			want: "34.5K",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateQuantity(tt.args.quantity); got != tt.want {
				t.Errorf("TruncateQuantity() = %v, want %v", got, tt.want)
			}
		})
	}
}
