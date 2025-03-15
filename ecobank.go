package ecobank

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	defaultBaseURL = "https://developer.ecobank.com/corporateapi/"
	userAgent      = "go-ecobank"
	origin         = "developer.ecobank.com"
	contentType    = "application/json"
)

// Client manages communication with the Ecobank API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *retryablehttp.Client

	// Base URL for API requests.
	baseURL *url.URL

	// Token for authenticating requests.
	token          string
	tokenExpiresAt time.Time

	// Credentials for requesting a token.
	username, password, labKey string

	// UserAgent is set in the User-Agent header of all requests.
	UserAgent string

	AccessToken *AccessTokenService
	Account     *AccountService
	Payment     PaymentService
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

	if err := c.setBaseURL(defaultBaseURL); err != nil {
		return nil, err
	}

	c.AccessToken = &AccessTokenService{client: c}
	c.Account = &AccountService{client: c}
	c.Payment = PaymentService{client: c}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// DoRequest sends an API request and returns the API response.
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

	token, resp, err := c.AccessToken.GetToken(ctx, req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	c.token = token.Token
	c.tokenExpiresAt = time.Now().Add(1 * time.Hour).Add(50 * time.Minute)

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

		log.Printf("Request body: %s", body)
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

// Do sends an API request and returns the API response.
func (c *Client) Do(req *retryablehttp.Request, v any) (*Response, error) {
	if c.token == "" || time.Now().After(c.tokenExpiresAt) {
		if c.username == "" && c.password == "" {
			return nil, errors.New("no token provided")
		}
		if err := c.Login(req.Context()); err != nil {
			return nil, err
		}
	}
	req.Header.Add("Authorization", "Bearer "+c.token)

	resp, err := c.doRequest(req, v)
	if err != nil {
		return nil, err
	}

	// TODO: Handle rate limiting
	// TODO: Handle token expiration, and reauthenticate

	return resp, nil
}

func (c *Client) doRequest(req *retryablehttp.Request, v any) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// TODO: check response for errors
	//if statusOk := resp.StatusCode >= 200 && resp.StatusCode < 300; !statusOk {
	//	return resp, errors.New(resp.Status)
	//}

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
			var respData ResponseData
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

func (c *Client) ensureSecureHash(opt any) {
	if sh, ok := opt.(secureHasher); ok && sh.GetHash() == "" {
		sh.SetHash(generateSecureHashFrom(opt, c.labKey))

		if sh.GetHash() != "398d4f285cc33e12f035da19fa9d954be35afaf66816531c4f1a1aedd3c6f132a85c62b23ca12d7b9a99bf5a84fc69b66738289a70e8f8115e90ffaa060f4026" {
			log.Println("Secure hash does not match")
		}
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

	log.Println("String to hash:", b.String()+key)

	return generateSecureHash(b.String(), key)
}

// generateSecureHash generates a secure hash for the given data.
func generateSecureHash(data, key string) string {
	hash := sha512.Sum512([]byte(data + key))
	return hex.EncodeToString(hash[:])
}

// ResponseData represents the data returned by the API.
type ResponseData struct {
	ResponseCode    int             `json:"response_code"`
	ResponseMessage string          `json:"response_message"`
	ResponseContent json.RawMessage `json:"response_content"`
	ResponseTime    Time            `json:"response_timestamp"`
	Errors          ResponseError   `json:"errors"`
}

var emptyResponseContent = []byte{0x22, 0x22}

func unmarshalResponse(resp any, data *ResponseData) error {
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
