package slack

import (
	"bytes"
	"io"
	"net/http"
)

type mockHttp struct{}

func (c *mockHttp) Do(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader([]byte("{\"ok\":true}"))),
		StatusCode: http.StatusOK,
	}, nil
}
