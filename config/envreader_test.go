package config

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Foo string `mapstructure:"foo"`
}

func Test_env2Map(t *testing.T) {
	getEnv = func(key string) string {
		return "bar"
	}
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *map[string]interface{}
		wantErr bool
	}{
		{
			"decode mapstructure tags",
			args{v: TestStruct{Foo: "bar"}},
			&map[string]interface{}{
				"foo": "bar",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env2Map(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("env2Map() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("env2Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
