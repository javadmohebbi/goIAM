package cmds

// Package cmds contains command-line tools and HTTP helpers for interacting with the goIAM API.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

// const DefaultUserAgent = "goIAM-CLI/1.0 ({os}; {arch}) Go-http-client/2.0"

func BuildUserAgent() string {
	return fmt.Sprintf("goIAM-CLI/1.0 (%s; %s) Go-http-client/2.0", runtime.GOOS, runtime.GOARCH)
}

// request sends an HTTP request to the goIAM API with support for multiple methods.
//
// Parameters:
//   - method: the HTTP method (e.g., "POST", "PATCH", "GET").
//   - apiURL: pointer to the base API URL (e.g., https://api.example.com).
//   - path: endpoint path (e.g., /auth/register).
//   - data: a map of key-value pairs to encode as the request body (ignored for GET).
//   - token: optional Bearer token for Authorization header.
//   - headers: optional variadic map of custom headers (only the first map is used).
//
// Returns the *http.Response and error (if any).
func request(method string, apiURL *string, path string, data map[string]any, token string, headers ...map[string]string) (resp *http.Response, err error) {
	var body *bytes.Buffer
	if data != nil && (method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	} else {
		body = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, *apiURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", BuildUserAgent())
	if data != nil && (method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut) {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set(k, v)
		}
	}

	resp, err = http.DefaultClient.Do(req)
	return resp, err
}
