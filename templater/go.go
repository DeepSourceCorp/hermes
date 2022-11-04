package templater

import (
	"bytes"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type Go struct{}

func (*Go) Execute(pattern string, params interface{}) ([]byte, error) {
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"concatenateWords": ConcatenateWords,
		"duration":         Duration,
		"plural":           Plural,
		"pluralWord":       PluralWord,
		"truncateQuantity": TruncateQuantity,
		"escapeSlackText":  EscapeSlackText,
	}).Parse(pattern)

	if err != nil {
		log.Errorf("Failed to parse template %s pattern: %v", tmpl.Name(), err)
		return nil, err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, params)
	if err != nil {
		log.Errorf("Failed to execute template %s: %v", tmpl.Name(), err)
		return nil, err
	}

	return b.Bytes(), nil
}
