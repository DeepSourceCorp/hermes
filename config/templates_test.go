package config

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/deepsourcelabs/hermes/domain"
	"gopkg.in/yaml.v3"
)

var template1 = Template{
	ID:                 "id1",
	Path:               "./file1.txt",
	SupportedProviders: []domain.ProviderType{},
}

var template2 = Template{
	ID:                 "id2",
	Path:               "./file2.txt",
	SupportedProviders: []domain.ProviderType{},
}

var tConfig = TemplateConfig{
	[]Template{template1, template2},
}

var errTest = errors.New("test error")

func TestTemplateConfig_Validate(t *testing.T) {
	t.Run("template files exists [mocked]", func(t *testing.T) {
		osStat = func(_ string) (os.FileInfo, error) { return nil, nil }
		if err := tConfig.Validate(); err != nil {
			t.Errorf("TemplateConfig.Validate() unexpected error = %v,", err)
		}
	})

	t.Run("template file does not exist [mocked]", func(t *testing.T) {
		osStat = func(_ string) (os.FileInfo, error) { return nil, errTest }
		if err := tConfig.Validate(); err == nil {
			t.Errorf("TemplateConfig.Validate() unexpected error = %v,", err)
		}
	})
}

func TestTemplateConfig_ReadYAML(t *testing.T) {
	t.Run("template config read", func(t *testing.T) {
		osReadFile = func(_ string) ([]byte, error) {
			return yaml.Marshal(&tConfig)
		}
		got := TemplateConfig{}
		if err := got.ReadYAML("template_path"); err != nil {
			t.Errorf("TemplateConfig.ReadYAML() unexpected error = %v", err)
		}

		if !reflect.DeepEqual(got, tConfig) {
			t.Errorf("TemplateConfig.ReadYAML() expected = %v, got = %v", tConfig, got)
		}
	})

	t.Run("template config read error", func(t *testing.T) {
		osReadFile = func(_ string) ([]byte, error) { return nil, errTest }
		got := TemplateConfig{}
		if err := got.ReadYAML("test"); err == nil {
			t.Errorf("TemplateConfig.ReadYAML() unexpected error = %v,", err)
		}
	})
}

func TestInitTemplateConfig(t *testing.T) {
	t.Run("template config read all files exist [mocked]", func(t *testing.T) {
		osStat = func(_ string) (os.FileInfo, error) { return nil, nil }
		osReadFile = func(_ string) ([]byte, error) {
			return yaml.Marshal(&tConfig)
		}
		if err := InitTemplateConfig("test"); err != nil {
			t.Errorf("InitTemplateConfig() unexpected error = %v", err)
		}
	})

	t.Run("template config read fail all files exist [mocked]", func(t *testing.T) {
		osStat = func(_ string) (os.FileInfo, error) { return nil, nil }
		osReadFile = func(_ string) ([]byte, error) { return nil, errTest }
		if err := InitTemplateConfig("test"); err == nil {
			t.Error("InitTemplateConfig() expected error, got nil")
		}
	})
}
