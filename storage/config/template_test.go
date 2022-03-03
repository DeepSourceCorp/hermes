package config

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/deepsourcelabs/hermes/config"
	"github.com/deepsourcelabs/hermes/domain"
)

func Test_templateStore_GetByID(t *testing.T) {
	type fields struct {
		cfg        *config.TemplateCfg
		fnReadFile func(filename string) ([]byte, error)
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Template
		wantErr bool
	}{
		{
			name: "load template from config",
			fields: fields{
				cfg: &config.TemplateCfg{Templates: []config.Template{
					{
						ID:                 "abc",
						Path:               "xyz",
						Type:               domain.TemplateType("TTT"),
						SupportedProviders: []domain.ProviderType{domain.ProviderType("PPP")},
					},
				}},
				fnReadFile: func(filename string) ([]byte, error) {
					return []byte(filename), nil
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "abc",
			},
			want: &domain.Template{
				ID:                 "abc",
				Pattern:            "xyz",
				Type:               domain.TemplateType("TTT"),
				SupportedProviders: []domain.ProviderType{domain.ProviderType("PPP")},
			},
			wantErr: false,
		},
		{
			name: "read failed",
			fields: fields{
				cfg: &config.TemplateCfg{Templates: []config.Template{
					{
						ID:                 "abc",
						Path:               "xyz",
						Type:               domain.TemplateType("TTT"),
						SupportedProviders: []domain.ProviderType{domain.ProviderType("PPP")},
					},
				}},
				fnReadFile: func(filename string) ([]byte, error) {
					return nil, errors.New("test")
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "abc",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &templateStore{
				cfg:        tt.fields.cfg,
				fnReadFile: tt.fields.fnReadFile,
			}
			got, err := store.GetByID(tt.args.ctx, tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("templateStore.GetByID() got = %v, want %v", got, tt.want)
			}
			if err != nil != tt.wantErr {
				t.Errorf("templateStore.GetByID() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
