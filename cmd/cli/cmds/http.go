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

// post sends a JSON-encoded POST request to the goIAM API and prints the response.
//
// Parameters:
//   - apiURL: pointer to the base API URL (e.g., https://api.example.com)
//   - path: endpoint path (e.g., /auth/register)
//   - data: a map of key-value pairs to encode as the request body
//   - token: optional Bearer token to include in the Authorization header
//   - headers: optional variadic map of additional headers to attach to the request
//
// It prints the response body to stdout and reports any errors encountered.
func post(apiURL *string, path string, data map[string]any, token string, headers ...map[string]string) (resp *http.Response, err error) {
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", *apiURL+path, bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", BuildUserAgent())

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// it will overwrite the headers if sent from the caller
	// like authorization as well
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set(k, v)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return nil, err
	}
	// defer res.Body.Close()

	// result, _ := io.ReadAll(res.Body)
	// fmt.Println(string(result))

	return res, err
}
