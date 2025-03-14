package ecobank

// ClientOptionFunc is a function that configures a Client.
type ClientOptionFunc func(*Client) error

// WithBaseURL sets the base URL for API requests to a custom endpoint.
func WithBaseURL(urlStr string) ClientOptionFunc {
	return func(c *Client) error {
		return c.setBaseURL(urlStr)
	}
}

func WithToken(token string) ClientOptionFunc {
	return func(c *Client) error {
		c.token = token

		return nil
	}
}
