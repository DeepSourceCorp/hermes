package templater

import (
	"bytes"
	"text/template"
)

type Opts struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

type ITemplater interface {
	Parse(filename string) error
	Execute(params interface{}) (string, error)
}

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

func GetTemplater(opts *Opts) ITemplater {
	switch opts.Type {
	case "text":
		return new(textTemplater)

	}
	return nil
}
