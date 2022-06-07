package config

import (
	"os"
	"path"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var templateConfig *TemplateConfig

var (
	osReadFile readFileFn = os.ReadFile
	osStat     statFn     = os.Stat
)

type Template struct {
	ID                 string                `yaml:"id,omitempty"`
	Path               string                `yaml:"path,omitempty"`
	Type               domain.TemplateType   `yaml:"type,omitempty"`
	SupportedProviders []domain.ProviderType `yaml:"supported_providers"`
}

type TemplateConfig struct {
	Templates []Template `yaml:"templates"`
}

func (tc *TemplateConfig) Validate() error {
	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, t := range tc.Templates {
		_, err := osStat(path.Join(workingDir, t.Path))
		if err != nil {
			return err
		}
	}
	return nil
}

func (tc *TemplateConfig) ReadYAML(configPath string) error {
	configBytes, err := osReadFile(path.Join(configPath, "./template.yaml"))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configBytes, &tc)
}

func InitTemplateConfig(templateConfigPath string) error {
	tempConfig := new(TemplateConfig)
	if err := tempConfig.ReadYAML(templateConfigPath); err != nil {
		return err
	}
	if err := tempConfig.Validate(); err != nil {
		return err
	}
	templateConfig = tempConfig
	log.Info("loaded new template config")
	return nil
}

type TemplateConfigFactory interface {
	GetTemplateConfig() *TemplateConfig
}

type templateConfigFactory struct{}

func NewTemplateConfigFactory() TemplateConfigFactory {
	return &templateConfigFactory{}
}

func (*templateConfigFactory) GetTemplateConfig() *TemplateConfig {
	return templateConfig
}

func StartTemplateConfigWatcher(configPath string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := InitTemplateConfig(configPath); err != nil {
						log.Error("failed to reload config, ", err.Error())
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error(err.Error())
			}
		}
	}()

	if err := watcher.Add(configPath); err != nil {
		return err
	}

	<-done
	return nil
}
