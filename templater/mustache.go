package templater

import "github.com/cbroglie/mustache"

type mustacheTemplater struct {
	template *mustache.Template
}

func (tmplr *mustacheTemplater) Parse(filename string) error {
	t, err := mustache.ParseFile(filename)
	tmplr.template = t
	if err != nil {
		return err
	}
	return nil
}

func (tmplr *mustacheTemplater) Execute(params interface{}) (string, error) {
	return tmplr.template.Render(params)
}
