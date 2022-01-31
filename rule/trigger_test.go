package rule

import (
	"reflect"
	"testing"
)

func TestRuleTrigger_Evaluate(t *testing.T) {
	type fields struct {
		RuleExpression string
		accessors      []string
	}
	type args struct {
		payload map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			"should evaluate rule",
			fields{
				RuleExpression: `[a.b] == 1 && [x.y] == "x"`,
			},
			args{
				map[string]interface{}{
					"a.b": 1,
					"x.y": "x",
				},
			},
			true,
			false,
		},
		{
			"should be falsy when rule fails",
			fields{
				RuleExpression: `[a.b] == 1 && [x.y] == "x"`,
			},
			args{
				map[string]interface{}{
					"a.b": 2,
					"x.y": "x",
				},
			},
			false,
			false,
		},
		{
			"should error out when expression evaluates to non-bool",
			fields{
				RuleExpression: `1*2`,
			},
			args{
				map[string]interface{}{
					"a.b": 2,
					"x.y": "x",
				},
			},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RuleTrigger{
				RuleExpression: tt.fields.RuleExpression,
				accessors:      tt.fields.accessors,
			}
			got, err := c.Evaluate(tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("RuleTrigger.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RuleTrigger.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleTrigger_extractAccessors(t *testing.T) {
	type fields struct {
		RuleExpression string
		accessors      []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			"must extract accessors correctly",
			fields{RuleExpression: `[a.b] == "foo" && [x.y.z] == "bar"`},
			[]string{
				"a.b",
				"x.y.z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RuleTrigger{
				RuleExpression: tt.fields.RuleExpression,
				accessors:      tt.fields.accessors,
			}
			if got := c.extractAccessors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RuleTrigger.extractAccessors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuleTrigger_MakeParams(t *testing.T) {
	type fields struct {
		RuleExpression string
		accessors      []string
	}
	type args struct {
		eventJSON []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			"should generate the correct parameter map",
			fields{
				accessors: []string{"a.b", "x.y.z"},
			},
			args{
				[]byte(`{"a":{"b":1}, "x":{"y":{"z":"foo"}}}`),
			},
			map[string]interface{}{
				"a.b":   1.0,
				"x.y.z": "foo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RuleTrigger{
				RuleExpression: tt.fields.RuleExpression,
				accessors:      tt.fields.accessors,
			}
			if got := c.MakeParams(tt.args.eventJSON); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RuleTrigger.MakeParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
