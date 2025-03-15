package ecobank

import (
	"context"
	"net/http"
)

// StatusService handles communication with the status related methods of the Ecobank API.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#88e33130-bba2-4bc4-9d77-e012b7124911
type StatusService struct {
	client *Client
}

// TransactionStatus represents the status of a transaction.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#758a9aef-edc6-45de-8ab0-1631c80936a1
type TransactionStatus struct {
	RequestType      string `json:"requestType"`
	AffiliateCode    string `json:"affiliateCode"`
	Amount           int    `json:"amount"`
	Currency         string `json:"currency"`
	Status           string `json:"status"`
	StatusCode       string `json:"statusCode"`
	StatusReason     string `json:"statusReason"`
	TransactionRefNo string `json:"transactionRefNo"`
}

// StatusOptions specifies the request parameters to get the status of a transaction.
type StatusOptions struct {
	ClientID  string `json:"clientId"`
	RequestID string `json:"requestId"`

	secureHashOption
}

// GetTransactionStatus gets the status of a transaction.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#758a9aef-edc6-45de-8ab0-1631c80936a1
func (s *StatusService) GetTransactionStatus(ctx context.Context, opt *StatusOptions) (*TransactionStatus, *Response, error) {
	return DoRequest[TransactionStatus](ctx, s.client, http.MethodPost, "merchant/txns/status", opt)
}

// ETokenStatusOptions specifies the request parameters to get the status of a token.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#5f689c50-1c6a-4c47-af68-83ac66d8315f
type ETokenStatusOptions struct {
	RequestID     string `json:"requestId"`
	AffiliateCode string `json:"affiliateCode"`

	secureHashOption
}

// GetETokenStatus gets the status of a token.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#5f689c50-1c6a-4c47-af68-83ac66d8315f
func (s *StatusService) GetETokenStatus(ctx context.Context, opt *ETokenStatusOptions) (*string, *Response, error) {
	return DoRequest[string](ctx, s.client, http.MethodPost, "merchant/etoken/status", opt)
}
