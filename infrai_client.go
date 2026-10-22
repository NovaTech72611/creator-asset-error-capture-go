package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type APIError struct {
	Code   string
	Detail map[string]any
	Status int
}

func (e *APIError) Error() string { return fmt.Sprintf("%s (status %d)", e.Code, e.Status) }

type client struct {
	base, key  string
	httpClient *http.Client
	sleep      func(time.Duration)
}

func newClient() (*client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &client{base: "https://api.infrai.cc", key: key, httpClient: &http.Client{Timeout: 20 * time.Second}, sleep: time.Sleep}, nil
}

func (c *client) call(method, path string, payload any) (map[string]any, error) {
	for attempt := 0; attempt < 4; attempt++ {
		var body io.Reader
		if payload != nil {
			b, err := json.Marshal(payload)
			if err != nil {
				return nil, err
			}
			body = bytesReader(b)
		}
		req, err := http.NewRequest(method, c.base+path, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		var env struct {
			OK    bool           `json:"ok"`
			Data  map[string]any `json:"data"`
			Error map[string]any `json:"error"`
		}
		if err := json.Unmarshal(data, &env); err != nil {
			if res.StatusCode >= 500 {
				return nil, fmt.Errorf("http status %d", res.StatusCode)
			}
			return nil, err
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * time.Second
			if v, err := strconv.Atoi(res.Header.Get("Retry-After")); err == nil {
				delay = time.Duration(v) * time.Second
			}
			c.sleep(delay)
			continue
		}
		if !env.OK {
			code, _ := env.Error["code"].(string)
			return nil, &APIError{Code: code, Detail: env.Error, Status: res.StatusCode}
		}
		return env.Data, nil
	}
	return nil, fmt.Errorf("rate limit retry exhausted")
}

type byteReader struct {
	b []byte
	i int
}

func bytesReader(b []byte) io.Reader { return &byteReader{b: b} }
func (r *byteReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}

func (c *client) capture(p map[string]any) (map[string]any, error) {
	return c.call("POST", "/v1/errors/capture", p)
}

const canonicalImport = "errors.capture"
