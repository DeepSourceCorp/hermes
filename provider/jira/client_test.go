package jira

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"testing"

// 	"github.com/deepsourcelabs/hermes/provider"
// )

// type mockHttp struct{}

// func (c *mockHttp) Do(request *http.Request) (*http.Response, error) {
// 	return &http.Response{
// 		Body:       io.NopCloser(bytes.NewReader([]byte("{\"ok\":true}"))),
// 		StatusCode: http.StatusOK,
// 	}, nil
// }

// func Test_client_SendMessage(t *testing.T) {
// 	type tfields struct {
// 		httpClient provider.IHTTPClient
// 	}
// 	type args struct {
// 		request *postIssueRequest
// 	}
// 	tests := []struct {
// 		name    string
// 		tfields tfields
// 		args    args
// 		want    interface{}
// 		wantErr bool
// 	}{
// 		{
// 			name: "should trigger http client with appropriate params",
// 			args: args{&postIssueRequest{
// 				BearerToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik16bERNemsxTVRoRlFVRTJRa0ZGT0VGRk9URkJOREJDTVRRek5EZzJSRVpDT1VKRFJrVXdNZyJ9.eyJodHRwczovL2F0bGFzc2lhbi5jb20vb2F1dGhDbGllbnRJZCI6IkdmVnNoTkJjNUhndjZQRWhRVjNDVGJNcnFBa1RJNXN0IiwiaHR0cHM6Ly9hdGxhc3NpYW4uY29tL2VtYWlsRG9tYWluIjoiZGVlcHNvdXJjZS5pbyIsImh0dHBzOi8vYXRsYXNzaWFuLmNvbS9zeXN0ZW1BY2NvdW50SWQiOiI2MjFjNjZlZmRiNThjMTAwNjg3N2UyYTIiLCJodHRwczovL2F0bGFzc2lhbi5jb20vc3lzdGVtQWNjb3VudEVtYWlsRG9tYWluIjoiY29ubmVjdC5hdGxhc3NpYW4uY29tIiwiaHR0cHM6Ly9hdGxhc3NpYW4uY29tL3ZlcmlmaWVkIjp0cnVlLCJodHRwczovL2F0bGFzc2lhbi5jb20vZmlyc3RQYXJ0eSI6ZmFsc2UsImh0dHBzOi8vYXRsYXNzaWFuLmNvbS8zbG8iOnRydWUsImlzcyI6Imh0dHBzOi8vYXRsYXNzaWFuLWFjY291bnQtcHJvZC5wdXMyLmF1dGgwLmNvbS8iLCJzdWIiOiJhdXRoMHw1ZjMyNGI1NThkODllMzAwNDYxODRhNjgiLCJhdWQiOiJhcGkuYXRsYXNzaWFuLmNvbSIsImlhdCI6MTY0NjA1MzU0MCwiZXhwIjoxNjQ2MDU3MTQwLCJhenAiOiJHZlZzaE5CYzVIZ3Y2UEVoUVYzQ1RiTXJxQWtUSTVzdCIsInNjb3BlIjoicmVhZDppc3N1ZTpqaXJhLXNvZnR3YXJlIHJlYWQ6aXNzdWU6amlyYSByZWFkOnByb2plY3QucHJvcGVydHk6amlyYSByZWFkOnByb2plY3QuZmVhdHVyZTpqaXJhIHdyaXRlOmF0dGFjaG1lbnQ6amlyYSB3cml0ZTpjb21tZW50OmppcmEgd3JpdGU6Y29tbWVudC5wcm9wZXJ0eTpqaXJhIHdyaXRlOmlzc3VlOmppcmEifQ.jLlYHFkOZuX6FIv-mwXbUO7bcurWRwQLRdl2yvsQ0N0qooE2pmRm0MuO42VawJ47HyNO_w5LqXRvAvLDqs0XszipFXwSzQzItnYTpiGsFVG7osoQNaWP-kFLEr5flOYnVzk4xnp_8q_hLdB-NuHc9xi1AMvbgIaJCbKBqjgRK8dQsCx9NVQX66wMoV9Dvm5kJ3jhs1y022DVqy4d2TAk7aFEHeKUgd1HZBnZUVkbpgVHjVl_c3lUV6mXcDUSzg3nnrY_qdQT_lmN8uhRYogYSkPkuD1Wd9lHWXamQ8K1PQIGJwLspbptuTIwzwo55X93TpzJUSTpUC8eyzf39gKghA",
// 				CloudID:     "09322103-001d-432e-b0e6-96d44dd53262",
// 				Fields: fields{
// 					Project:   project{Key: "TEST"},
// 					IssueType: issueType{Name: "Bug"},
// 					Summary:   "--summary--",
// 					Description: map[string]interface{}{
// 						"version": 1,
// 						"type":    "doc",
// 						"content": []map[string]interface{}{
// 							{
// 								"type": "paragraph",
// 								"content": []map[string]interface{}{
// 									{
// 										"type": "text",
// 										"text": "Hello",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			}},
// 			tfields: tfields{
// 				httpClient: http.DefaultClient,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			payload, _ := json.Marshal(tt.args.request)
// 			fmt.Println(string(payload))
// 			c := &client{
// 				httpClient: tt.tfields.httpClient,
// 			}
// 			_, err := c.SendMessage(tt.args.request)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("client.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }
