package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"bizzmod-cli/internal/config"
)

const cliUserAgent = "uproc-cli/1.0"

type Client struct {
	httpClient *http.Client
	config     config.Config
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		config:     cfg,
	}
}

func (c *Client) Do(method, path string, body []byte) ([]byte, int, error) {
	url := strings.TrimRight(c.config.APIURL, "/") + path

	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Authorization", "Bearer "+c.config.UserAPIKey)
	req.Header.Set("x-client-app", "uproc-cli")
	req.Header.Set("User-Agent", cliUserAgent)
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	if res.StatusCode >= 400 {
		return respBody, res.StatusCode, fmt.Errorf("http %d", res.StatusCode)
	}

	return respBody, res.StatusCode, nil
}

func (c *Client) DoMultipart(method, path string, contentType string, reader io.Reader) ([]byte, int, error) {
	url := strings.TrimRight(c.config.APIURL, "/") + path

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Authorization", "Bearer "+c.config.UserAPIKey)
	req.Header.Set("x-client-app", "uproc-cli")
	req.Header.Set("User-Agent", cliUserAgent)
	req.Header.Set("Content-Type", contentType)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	if res.StatusCode >= 400 {
		return respBody, res.StatusCode, fmt.Errorf("http %d", res.StatusCode)
	}

	return respBody, res.StatusCode, nil
}

func (c *Client) BaseURL() string {
	return strings.TrimRight(c.config.APIURL, "/")
}
