package domain

import (
	"io/ioutil"
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

func TestExecuteJinja(t *testing.T) {
	// a mock type for testing templates
	type result struct {
		File        string
		Issue       string
		Description string
	}

	// prepare params
	params := map[string]interface{}{
		"title": "Hermes",
		"results": []*result{
			{
				File:        "demo.py",
				Issue:       "PY-0014",
				Description: "Demo issue",
			},
		},
	}

	tests := []struct {
		name         string
		templateType string
		template     string
		want         string
		payload      map[string]interface{}
	}{
		{"jinja_must work with pattern and data", "jinja", "./testdata/template.txt", "./testdata/template_executed.txt", map[string]interface{}{"deepsource": params}},
		{"jinja_must work with pattern and no data", "jinja", "./testdata/template.txt", "./testdata/template_executed_nil.txt", map[string]interface{}{"deepsource": nil}},
		{"jinja_must work with pattern and no params", "jinja", "./testdata/template.txt", "./testdata/template_executed_nil.txt", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := &Template{
				Type: TemplateType(tt.templateType),
			}
			templater := temp.GetTemplater()

			// read pattern from file
			pattern, err := ioutil.ReadFile(tt.template)
			if err != nil {
				t.Errorf("couldn't open template: %v\n", err)
			}

			b, err := templater.Execute(string(pattern), tt.payload)
			assert.Nil(t, err)
			got := string(b)

			// read result from file
			res, err := ioutil.ReadFile(tt.want)
			if err != nil {
				t.Errorf("couldn't open result: %v\n", err)
			}
			want := string(res)

			if got != want {
				t.Errorf("got: %v, want: %v\n", got, want)
			}
		})
	}
}
