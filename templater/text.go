package templater

import (
	"bytes"
	"text/template"
)

type textTemplater struct {
	template *template.Template
}

func (tmplr *textTemplater) Parse(path string) error {
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	tmplr.template = t
	return nil
}

func (tmplr *textTemplater) Execute(params interface{}) (string, error) {
	buf := new(bytes.Buffer)
	tmplr.template.Execute(buf, params)
	return buf.String(), nil
}
