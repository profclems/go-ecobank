package ecobank

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	defaultBaseURL = "https://developer.ecobank.com/corporateapi/"
	userAgent      = "go-ecobank"
	origin         = "developer.ecobank.com"
	contentType    = "application/json"

	// by default, token expires in 2 hours
	defaultTokenExpiry = 7200 * time.Second
)

// Client manages communication with the Ecobank API.
type Client struct {
	// HTTP client used to communicate with the API.
	client         *retryablehttp.Client
	disableRetries bool

	// Base URL for API requests.
	baseURL *url.URL

	// Token for authenticating requests.
	tokenMu        sync.RWMutex
	token          string
	tokenExpiresAt time.Time

	// Credentials for requesting a token.
	username, password, labKey string

	// UserAgent is set in the User-Agent header of all requests.
	UserAgent string

	Auth       *AuthService
	Account    *AccountService
	Payment    *PaymentService
	Remittance *RemittanceService
	Status     *StatusService
}

// getToken returns the token and expiry time.
// It is safe for concurrent access since it obtains a lock before reading.
func (c *Client) getToken() (string, time.Time) {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token, c.tokenExpiresAt
}

// setToken sets the token and expiry time.
// It is safe for concurrent access since it obtains a lock before writing.
func (c *Client) setToken(token string, expiresAt time.Time) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token, c.tokenExpiresAt = token, expiresAt
}

// NewClient returns a new Ecobank API client.
func NewClient(username, password, labKey string, opts ...ClientOptionFunc) (*Client, error) {
	c := &Client{
		username:  username,
		password:  password,
		labKey:    labKey,
		UserAgent: userAgent,
	}

	c.client = retryablehttp.NewClient()
	c.client.RetryWaitMin = 100 * time.Millisecond
	c.client.RetryWaitMax = 400 * time.Millisecond
	c.client.RetryMax = 5
	c.client.Logger = nil
	c.client.CheckRetry = c.retryHTTPCheck
	c.client.ErrorHandler = retryablehttp.PassthroughErrorHandler

	if err := c.setBaseURL(defaultBaseURL); err != nil {
		return nil, err
	}

	c.Auth = &AuthService{client: c}
	c.Account = &AccountService{client: c}
	c.Payment = &PaymentService{client: c}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// DoRequest sends an HTTP API request using the provided client and returns the response.
//
// This is a generic function designed to handle API requests dynamically while decoding
// the response into the expected type `T`. It abstracts away the request construction
// and execution process, making API calls more reusable and consistent.
//
// Parameters:
// - `ctx` (context.Context): The request context, allowing for request cancellation, timeouts, etc.
// - `client` (*Client): The API client used to send the request.
// - `method` (string): The HTTP method (e.g., "GET", "POST", "PUT", "DELETE").
// - `path` (string): The endpoint path for the API request.
// - `opt` (any): The request body or query parameters (can be nil).
//
// Returns:
// - `*T`: A pointer to the parsed response body, unmarshaled into the expected type `T`.
// - `*Response`: The full HTTP response, containing metadata such as headers and status code.
// - `error`: An error if the request fails or response decoding encounters an issue.
//
// Workflow:
// 1. Calls `client.NewRequest()` to create an HTTP request with the provided method, path, and options.
// 2. Executes the request using client.Do(), which sends the request and decodes the response.
// 3. If the request succeeds, the response body is unmarshaled into `T` and returned along with the HTTP response.
// 4. If any error occurs (e.g., network failure, non-200 status code, JSON decoding issues), it is returned.
//
// Example:
// ```go
//
//	type User struct {
//	    ID   int    `json:"id"`
//	    Name string `json:"name"`
//	}
//
// client := NewClient("https://api.example.com")
// ctx := context.Background()
//
// user, resp, err := DoRequest[User](ctx, client, "GET", "/users/123", nil)
//
//	if err != nil {
//	    log.Fatalf("Request failed: %v", err)
//	}
//
// fmt.Printf("User: %+v, Status: %d\n", user, resp.StatusCode)
// ```
func DoRequest[T any](ctx context.Context, client *Client, method, path string, opt any) (*T, *Response, error) {
	req, err := client.NewRequest(ctx, method, path, opt)
	if err != nil {
		return nil, nil, err
	}

	var respT T

	response, err := client.Do(req, &respT)
	if err != nil {
		return nil, response, err
	}

	return &respT, response, nil
}

// Login authenticates the client and stores the access token in the client.
func (c *Client) Login(ctx context.Context) error {
	req := &AccessTokenOptions{
		UserID:   c.username,
		Password: c.password,
	}

	token, resp, err := c.Auth.GetAccessToken(ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	expiry, err := getTokenExpiry(c.token)
	if err != nil {
		// set a default expiry time
		expiry = time.Now().Add(defaultTokenExpiry)
	}

	c.setToken(token.Token, expiry)

	return nil
}

// setBaseURL sets the base URL for API requests to a custom endpoint.
func (c *Client) setBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	// Update the base URL of the client.
	c.baseURL = baseURL

	return nil
}

// BaseURL returns a copy of the baseURL
func (c *Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
}

// NewRequest creates an API request.
func (c *Client) NewRequest(ctx context.Context, method, path string, opts any) (*retryablehttp.Request, error) {
	u := *c.baseURL

	unescaped, err := url.PathUnescape(path)
	if err != nil {
		return nil, err
	}

	u.RawPath = c.baseURL.Path + path
	u.Path = c.baseURL.Path + unescaped

	headers := make(http.Header)

	if c.UserAgent != "" {
		headers.Set("User-Agent", userAgent)
	}

	headers.Set("Content-Type", contentType)
	headers.Set("Accept", contentType)
	headers.Set("Origin", origin)

	var body any

	if opts != nil {
		c.ensureSecureHash(opts)
		body, err = json.Marshal(opts)
		if err != nil {
			return nil, err
		}
	}

	req, err := retryablehttp.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	for key, values := range headers {
		req.Header[key] = values
	}

	return req, nil
}

// Do sends an authenticated API request, ensuring a valid token is set and adding it to the request header.
//
// It processes the API response, unmarshaling the `response_content` field into v, while extracting other response
// metadata such as `response_code`, `response_message`, and `response_timestamp` into a *Response object.
// The API response is expected to follow the format:
//
//	{
//	  "response_code": "00",
//	  "response_message": "Success",
//	  "response_content": { ... },
//	  "response_timestamp": "2025-03-16T12:34:56Z"
//	}
//
// The `response_content` field is dynamically unmarshaled into v, while `response_code`, `response_message`,
// and `response_timestamp` are stored in the returned *Response alongside the underlying HTTP response.
//
// If the API response contains an `errors` field, it is returned as an error of type ResponseError.
// Use `errors.As(err, &ResponseError)` to extract the error details.
//
// Example:
//
//	var result SomeResponseType
//	resp, err := client.Do(req, &result)
//	if err != nil {
//		var apiErr ResponseError
//		if errors.As(err, &apiErr) {
//			log.Println("API error:", apiErr)
//		} else {
//			log.Println("Request failed:", err)
//		}
//	}
func (c *Client) Do(req *retryablehttp.Request, v any) (*Response, error) {
	token, expiry := c.getToken()
	// authenticate if token is not set or has expired
	if token == "" || (!expiry.IsZero() && time.Now().After(expiry)) {
		if c.username == "" && c.password == "" {
			return nil, errors.New("token expired")
		}
		if err := c.Login(req.Context()); err != nil {
			return nil, fmt.Errorf("failed to re-authenticate: %w", err)
		}

		token, _ = c.getToken()
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := c.doRequest(req, v)
	if err != nil {
		return nil, err
	}

	// TODO: Handle rate limiting

	return resp, nil
}

func (c *Client) doRequest(req *retryablehttp.Request, v any) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	r := newResponse(resp)

	if v != nil {
		defer func() {
			err = errors.Join(err, resp.Body.Close())
		}()
		defer func() {
			err = errors.Join(err, checkErr1(io.Copy(io.Discard, resp.Body)))
		}()

		if _, ok := v.(*BearerToken); ok {
			err = json.NewDecoder(resp.Body).Decode(v)
		} else {
			var respData responseData
			err = json.NewDecoder(resp.Body).Decode(&respData)
			if err == nil {
				r.Code = respData.ResponseCode
				r.Message = respData.ResponseMessage
				r.Time = respData.ResponseTime

				if respData.Errors != nil {
					return r, &respData.Errors
				}
				err = unmarshalResponse(v, &respData)
			}
		}
	}

	return r, err
}

// retryHTTPCheck provides a callback for Client.CheckRetry which
// will retry both rate limit (429) and server (>= 500) errors.
func (c *Client) retryHTTPCheck(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		return false, err
	}
	if !c.disableRetries && (resp.StatusCode == 429 || resp.StatusCode >= 500) {
		return true, nil
	}
	return false, nil
}

func (c *Client) ensureSecureHash(opt any) {
	if sh, ok := opt.(secureHasher); ok && sh.GetHash() == "" {
		sh.SetHash(generateSecureHashFrom(opt, c.labKey))
	}
}

// generateSecureHashFrom generates a secure hash for the given struct.
func generateSecureHashFrom(v any, key string) string {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	var b strings.Builder

	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)
		// check if it's a struct and has the name PaymentHeader
		// For payment, the secure hash is generated from the PaymentHeader struct
		if typ.Kind() == reflect.Struct && fieldType.Tag.Get("json") == "paymentHeader" {
			return generateSecureHashFrom(fieldValue.Interface(), key)
		}

		// skip unexported fields, anonymous fields, fields with securehash tag set to ignore, and fields with json tag set to "-"
		if !fieldValue.CanInterface() ||
			fieldType.Anonymous ||
			fieldType.Tag.Get("securehash") == "ignore" ||
			fieldType.Tag.Get("json") == "-" ||
			fieldType.Tag.Get("json") == "secureHash" {
			continue
		}

		b.WriteString(formatToStr(fieldValue.Interface()))
	}

	return generateSecureHash(b.String(), key)
}

// generateSecureHash generates a secure hash for the given data.
func generateSecureHash(data, key string) string {
	hash := sha512.Sum512([]byte(data + key))
	return hex.EncodeToString(hash[:])
}

type responseData struct {
	ResponseCode    int             `json:"response_code"`
	ResponseMessage string          `json:"response_message"`
	ResponseContent json.RawMessage `json:"response_content"`
	ResponseTime    Time            `json:"response_timestamp"`
	Errors          ResponseError   `json:"errors"`
}

var emptyResponseContent = []byte{0x22, 0x22}

func unmarshalResponse(resp any, data *responseData) error {
	if data.ResponseContent == nil || bytes.Equal(data.ResponseContent, emptyResponseContent) {
		return nil
	}
	return json.Unmarshal(data.ResponseContent, resp)
}

type Response struct {
	*http.Response

	// Code is the response_code returned by the API as part of the response payload but not the HTTP status code.
	// 000 or 200 means successful submission and any other code means an error occurred.
	Code int
	// Message is the response_message returned by the API as part of the response payload but not the HTTP status message.
	Message string
	// Time is the response_timestamp returned by the API as part of the response payload.
	Time Time
}

func newResponse(r *http.Response) *Response {
	return &Response{
		Response: r,
	}
}
