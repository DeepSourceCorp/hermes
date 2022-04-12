package templater

import "github.com/flosch/pongo2"

type Jinja struct{}

func (*Jinja) Execute(pattern string, params interface{}) ([]byte, error) {
	tmpl, err := pongo2.FromString(pattern)
	if err != nil {
		return nil, err
	}

	out, err := tmpl.Execute(params.(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	return []byte(out), nil
}
