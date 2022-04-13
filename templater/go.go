package templater

import (
	"bytes"
	"html/template"
)

type Go struct{}

func (*Go) Execute(pattern string, params interface{}) ([]byte, error) {
	tmpl, err := template.New("template").Parse(pattern)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, params)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
