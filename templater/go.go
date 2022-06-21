package templater

import (
	"bytes"
	"text/template"
)

type Go struct{}

func (*Go) Execute(pattern string, params interface{}) ([]byte, error) {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"concatenateWords": ConcatenateWords,
		"duration":         Duration,
		"plural":           Plural,
		"pluralWord":       PluralWord,
		"truncateQuantity": TruncateQuantity,
	}).Parse(pattern)
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
