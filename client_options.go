package ecobank

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// ClientOptionFunc is a function that configures a Client.
type ClientOptionFunc func(*Client) error

// WithBaseURL sets the base URL for API requests to a custom endpoint.
func WithBaseURL(urlStr string) ClientOptionFunc {
	return func(c *Client) error {
		return c.setBaseURL(urlStr)
	}
}

// WithToken sets the token for the client.
// It also sets the token expiry time by decoding the token and extracting the expiry time.
func WithToken(token string) ClientOptionFunc {
	return func(c *Client) (err error) {
		c.token = token
		c.tokenExpiresAt, err = getTokenExpiry(token)
		return err
	}
}

// WithTokenAndExpiry sets the token and expiry time for the client.
func WithTokenAndExpiry(token string, expiresAt time.Time) ClientOptionFunc {
	return func(c *Client) error {
		c.token, c.tokenExpiresAt = token, expiresAt
		return nil
	}
}

// WithHTTPClient sets the HTTP client for the client.
func WithHTTPClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *Client) error {
		c.client.HTTPClient = httpClient
		return nil
	}
}

// WithRetryableClient sets the retryable client for the client.
func WithRetryableClient(client *retryablehttp.Client) ClientOptionFunc {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

// WithUserAgent sets the User-Agent header for the client.
func WithUserAgent(userAgent string) ClientOptionFunc {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithDisableRetries disables retries for the client.
func WithDisableRetries() ClientOptionFunc {
	return func(c *Client) error {
		c.disableRetries = true
		return nil
	}
}

// WithRetryPolicy sets the retry policy for the client.
func WithRetryPolicy(retry retryablehttp.CheckRetry) ClientOptionFunc {
	return func(c *Client) error {
		c.client.CheckRetry = retry
		return nil
	}
}

// WithBackoff sets the backoff strategy for the client.
func WithBackoff(backoff retryablehttp.Backoff) ClientOptionFunc {
	return func(c *Client) error {
		c.client.Backoff = backoff
		return nil
	}
}
