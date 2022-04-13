package templater

import "github.com/hoisie/mustache"

type Mustache struct{}

func (*Mustache) Execute(pattern string, params interface{}) ([]byte, error) {
	str := mustache.Render(pattern, params)
	return []byte(str), nil
}
