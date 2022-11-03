package templater

import "testing"

func TestEscapeSlackText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "escapes & to &amp;",
			args: args{text: "&"},
			want: "&amp;",
		},
		{
			name: "escapes < to &lt;",
			args: args{text: "<"},
			want: "&lt;",
		},
		{
			name: "escapes > to &gt;",
			args: args{text: ">"},
			want: "&gt;",
		},
		{
			name: "escapes multiple special characters",
			args: args{text: "ab&d<e>but&<z>"},
			want: "ab&amp;d&lt;e&gt;but&amp;&lt;z&gt;",
		},
		{
			name: "doesn't escape if no special characters are present",
			args: args{text: "this.has:no'slack'special[chars]"},
			want: "this.has:no'slack'special[chars]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeSlackText(tt.args.text); got != tt.want {
				t.Errorf("EscapeSlackText() = %v, want %v", got, tt.want)
			}
		})
	}
}
