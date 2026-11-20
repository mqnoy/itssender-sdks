package itssender

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultBaseURL is the default production gateway endpoint for itssender.
const DefaultBaseURL = "https://itssender.imzalab.com"

// Option configures an itssender API Client.
type Option func(*Client)

// WithAPIKey sets the API key used for X-API-Key authentication header.
func WithAPIKey(apiKey string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithBearerToken sets the JWT bearer token used for Authorization header.
func WithBearerToken(token string) Option {
	return func(c *Client) {
		c.bearerToken = token
	}
}

// WithBaseURL sets the base URL for API requests. Defaults to DefaultBaseURL ("https://itssender.imzalab.com").
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the request timeout on the HTTP client. Defaults to 30 seconds.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
		if c.httpClient != nil {
			c.httpClient.Timeout = timeout
		}
	}
}

// WithRetry sets the maximum number of retries for transient HTTP errors (429, 5xx). Defaults to 3.
func WithRetry(maxRetries int) Option {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// Client is the primary entrypoint for interacting with the itssender WhatsApp Gateway API.
type Client struct {
	baseURL     string
	apiKey      string
	bearerToken string
	httpClient  *http.Client
	timeout     time.Duration
	maxRetries  int

	// Sub-services
	Auth      *AuthService
	APIKeys   *APIKeysService
	Devices   *DevicesService
	Messages  *MessagesService
	Groups    *GroupsService
	Contacts  *ContactsService
	Campaigns *CampaignsService
	Analytics *AnalyticsService
	OTP       *OTPService
	Templates *TemplatesService
	Webhooks  *WebhooksService
}

// NewClient constructs an idiomatic itssender API client with configured options.
func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:    DefaultBaseURL,
		timeout:    30 * time.Second,
		maxRetries: 3,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: c.timeout,
		}
	} else if c.timeout > 0 {
		c.httpClient.Timeout = c.timeout
	}

	// Initialize sub-services
	c.Auth = &AuthService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.Devices = &DevicesService{client: c}
	c.Messages = &MessagesService{client: c}
	c.Groups = &GroupsService{client: c}
	c.Contacts = &ContactsService{client: c}
	c.Campaigns = &CampaignsService{client: c}
	c.Analytics = &AnalyticsService{client: c}
	c.OTP = &OTPService{client: c}
	c.Templates = &TemplatesService{client: c}
	c.Webhooks = &WebhooksService{client: c}

	return c
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// SetBearerToken dynamically updates the bearer token (e.g. after login or refresh).
func (c *Client) SetBearerToken(token string) {
	c.bearerToken = token
}

// SetAPIKey dynamically updates the API key.
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// do executes an HTTP request with automatic retries and JSON envelope decoding.
func (c *Client) do(ctx context.Context, method, endpoint string, body interface{}, out interface{}) error {
	fullURL := c.baseURL + "/" + strings.TrimLeft(endpoint, "/")

	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("itssender: failed to marshal request body: %w", err)
		}
	}

	attempts := c.maxRetries + 1
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	var resp *http.Response

	for attempt := 0; attempt < attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return fmt.Errorf("itssender: failed to construct request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "itssender-go/1.0.0")

		if c.apiKey != "" {
			req.Header.Set("X-API-Key", c.apiKey)
		}
		if c.bearerToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.bearerToken)
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < attempts-1 {
				c.backoffSleep(ctx, attempt)
				continue
			}
			return fmt.Errorf("itssender: network error: %w", lastErr)
		}

		// Check for retryable HTTP status codes: 429 Too Many Requests, 500, 502, 503, 504
		if resp.StatusCode == http.StatusTooManyRequests || (resp.StatusCode >= 500 && resp.StatusCode != http.StatusNotImplemented) {
			_ = resp.Body.Close()
			if attempt < attempts-1 {
				c.backoffSleep(ctx, attempt)
				continue
			}
			lastErr = fmt.Errorf("itssender: server returned status %d", resp.StatusCode)
			break
		}

		// Non-retryable status or success: proceed to decode response
		break
	}

	if resp == nil {
		if lastErr != nil {
			return lastErr
		}
		return errors.New("itssender: failed to execute request")
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("itssender: failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    string(respBytes),
		}

		// Try unmarshaling structured error
		var env struct {
			Success bool `json:"success"`
			Error   *struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
			Message string `json:"message"`
			Code    string `json:"code"`
		}

		if unmarshalErr := json.Unmarshal(respBytes, &env); unmarshalErr == nil {
			if env.Error != nil {
				apiErr.Code = env.Error.Code
				apiErr.Message = env.Error.Message
			} else if env.Message != "" {
				apiErr.Message = env.Message
				apiErr.Code = env.Code
			}
		}

		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}

		return apiErr
	}

	if out == nil || len(respBytes) == 0 {
		return nil
	}

	return decodeResponse(respBytes, out)
}

// decodeResponse extracts data from JSON response supporting envelope {"success":true,"data":...} and direct JSON.
func decodeResponse(respBytes []byte, out interface{}) error {
	var env struct {
		Success *bool           `json:"success"`
		Data    json.RawMessage `json:"data"`
		Error   *APIError       `json:"error"`
	}

	if err := json.Unmarshal(respBytes, &env); err == nil {
		if env.Success != nil && !*env.Success && env.Error != nil {
			return env.Error
		}
		if len(env.Data) > 0 && string(env.Data) != "null" {
			return json.Unmarshal(env.Data, out)
		}
	}

	// Fallback to direct unmarshal into target
	return json.Unmarshal(respBytes, out)
}

// backoffSleep sleeps with exponential backoff and jitter.
func (c *Client) backoffSleep(ctx context.Context, attempt int) {
	baseMs := int64(100 * (1 << attempt))
	if baseMs > 3000 {
		baseMs = 3000
	}
	jitterMs := int64(0)
	if n, err := rand.Int(rand.Reader, big.NewInt(50)); err == nil {
		jitterMs = n.Int64()
	}
	delay := time.Duration(baseMs+jitterMs) * time.Millisecond

	select {
	case <-time.After(delay):
	case <-ctx.Done():
	}
}

// buildURLWithQuery appends query parameters to an endpoint.
func buildURLWithQuery(endpoint string, query map[string]string) string {
	if len(query) == 0 {
		return endpoint
	}
	vals := url.Values{}
	for k, v := range query {
		if v != "" {
			vals.Add(k, v)
		}
	}
	encoded := vals.Encode()
	if encoded == "" {
		return endpoint
	}
	if strings.Contains(endpoint, "?") {
		return endpoint + "&" + encoded
	}
	return endpoint + "?" + encoded
}
