package templater

type Opts struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

type ITemplater interface {
	Parse(filename string) error
	Execute(params interface{}) (string, error)
}

// GetTemplater is a templater factory.
func GetTemplater(opts *Opts) ITemplater {
	switch opts.Type {
	case "text":
		return new(textTemplater)
	case "mustache":
		return new(mustacheTemplater)
	}
	return nil
}
