package ecobank

import "context"

type AuthService struct {
	client *Client
}

// AccessTokenOptions represents a request to get an access token.
type AccessTokenOptions struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
}

// BearerToken represents a response to an access token request.
type BearerToken struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// GetAccessToken gets an access token for the given user.
func (a *AuthService) GetAccessToken(ctx context.Context, opt *AccessTokenOptions) (*BearerToken, *Response, error) {
	req, err := a.client.NewRequest(ctx, "POST", "user/token", opt)
	if err != nil {
		return nil, nil, err
	}

	var token BearerToken

	resp, err := a.client.doRequest(req, &token)
	if err != nil {
		return nil, resp, err
	}

	return &token, resp, nil
}
