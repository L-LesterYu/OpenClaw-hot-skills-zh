package image

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newMockHTTPClient(fn roundTripFunc) *http.Client {
	return &http.Client{Transport: fn}
}

func jsonResponse(status int, body any) *http.Response {
	var data []byte
	switch value := body.(type) {
	case []byte:
		data = value
	case string:
		data = []byte(value)
	default:
		data, _ = json.Marshal(value)
	}

	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(data)),
	}
}
