package slack

import (
	"reflect"
	"testing"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
)

func TestClient_GetChannels(t *testing.T) {
	type fields struct {
		HTTPClient provider.IHTTPClient
	}
	type args struct {
		request *GetChannelsRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
		want1  domain.IError
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				HTTPClient: tt.fields.HTTPClient,
			}
			got, got1 := c.GetChannels(tt.args.request)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetChannels() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Client.GetChannels() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
