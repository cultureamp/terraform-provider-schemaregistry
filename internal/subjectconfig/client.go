package subjectconfig

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const contentType = "application/vnd.schemaregistry.v1+json"
const accept = "application/vnd.schemaregistry.v1+json, application/json"

// Client reads and writes subject-level config in Schema Registry.
type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
}

type configResponse struct {
	Normalize *bool `json:"normalize"`
}

type configRequest struct {
	Normalize bool `json:"normalize"`
}

// New creates a subject config client.
func New(baseURL string, httpClient *http.Client, username, password string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: httpClient,
		username:   username,
		password:   password,
	}
}

// GetNormalize returns the subject-level normalize value, or nil if the
// subject has no config or does not set it.
func (c *Client) GetNormalize(ctx context.Context, subject string) (*bool, error) {
	endpoint := fmt.Sprintf("%s/config/%s?defaultToGlobal=false", c.baseURL, url.PathEscape(subject))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", accept)
	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}

	var cfg configResponse
	if err := json.Unmarshal(body, &cfg); err != nil {
		return nil, err
	}
	return cfg.Normalize, nil
}

// SetNormalize sets the subject-level normalize value.
// Schema Registry merges config updates, so this keeps the subject's compatibility level.
func (c *Client) SetNormalize(ctx context.Context, subject string, normalize bool) error {
	payload, err := json.Marshal(configRequest{Normalize: normalize})
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/config/%s", c.baseURL, url.PathEscape(subject))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", accept)
	if c.username != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, string(body))
	}
	return nil
}
