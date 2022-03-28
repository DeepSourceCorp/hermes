package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name         string
		templateType string
		pattern      string
		payload      map[string]string
		want         string
	}{
		{"gotmpl_must work with pattern and data", "gotmpl", "Hi {{.name}}", map[string]string{"name": "Apollo"}, "Hi Apollo"},
		{"gotmpl_must work with pattern and no data", "gotmpl", "Hi {{.name}}", nil, "Hi <no value>"},
		{"mustache_must work with pattern and data", "mustache", "Hi {{name}}", map[string]string{"name": "Apollo"}, "Hi Apollo"},
		{"mustache_must work with pattern and no data", "mustache", "Hi {{name}}", nil, "Hi "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := &Template{
				Type: TemplateType(tt.templateType),
			}
			templater := temp.GetTemplater()

			b, err := templater.Execute(tt.pattern, tt.payload)
			assert.Nil(t, err)

			got := string(b)
			if got != tt.want {
				t.Errorf("got: %v, want: %v\n", got, tt.want)
			}
		})
	}
}
