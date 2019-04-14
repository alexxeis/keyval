package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Client is an API client
type Client struct {
	host       string
	userAgent  string
	httpClient *http.Client
}

// NewClient returns new API client
func NewClient(h string, ua string, c *http.Client) *Client {
	return &Client{
		host:       h,
		userAgent:  ua,
		httpClient: c,
	}
}

// newRequest creates new HTTP request
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.host+path, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

// process makes HTTP requests to API
func (c *Client) process(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if v == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(v)
}
