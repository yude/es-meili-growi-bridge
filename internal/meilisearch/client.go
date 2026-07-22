package meilisearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	host       string
	apiKey     string
}

func NewClient(host, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		host:       host,
		apiKey:     apiKey,
	}
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, c.host+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) doJSON(method, path string, reqBody, respBody interface{}) error {
	resp, err := c.doRequest(method, path, reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var meiliErr Error
		if err := json.NewDecoder(resp.Body).Decode(&meiliErr); err != nil {
			return fmt.Errorf("meilisearch error (status=%d)", resp.StatusCode)
		}
		return &meiliErr
	}

	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) doRaw(method, path string, reqBody interface{}) ([]byte, error) {
	resp, err := c.doRequest(method, path, reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var meiliErr Error
		if err := json.NewDecoder(resp.Body).Decode(&meiliErr); err != nil {
			return nil, fmt.Errorf("meilisearch error (status=%d)", resp.StatusCode)
		}
		return nil, &meiliErr
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) doNoBody(method, path string) ([]byte, error) {
	return c.doRaw(method, path, nil)
}
