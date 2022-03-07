package config

type AppConfig struct {
	// Server configuration
	Port int `mapstructure:"PORT"`

	TemplateConfigPath string `mapstructure:"TEMPLATE_CONFIG"`
}
