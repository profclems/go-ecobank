package ecobank

import (
	"context"
	"net/http"
)

// RemittanceService handles communication with the remittance related
// methods of the Ecobank API.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#acfe7f55-27aa-487d-ba1a-799ecb466bd7
type RemittanceService struct {
	client *Client
}

// Institution represents an Ecobank affiliate allowed to participate in cross-border transactions.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#eaeb6f0a-107d-4717-b202-b8eee1529b74
type Institution struct {
	InstitutionId   string `json:"institutionId"`
	InstitutionType string `json:"institutionType"`
	InstitutionName string `json:"institutionName"`
	CountryCode     string `json:"countryCode"`
}

// ListInstitutionsOptions specifies the request parameters to list institutions.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#eaeb6f0a-107d-4717-b202-b8eee1529b74
type ListInstitutionsOptions struct {
	RequestID          string `json:"requestId"`
	ClientID           string `json:"clientId"`
	AffiliateCode      string `json:"affiliateCode"`
	DestinationCountry string `json:"destinationCountry"`

	secureHashOption
}

// ListInstitutions returns the list of Ecobank affiliates allowed to participate in cross-border transactions.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#eaeb6f0a-107d-4717-b202-b8eee1529b74
func (s *RemittanceService) ListInstitutions(ctx context.Context, opt *ListInstitutionsOptions) ([]*Institution, *Response, error) {
	institutions, resp, err := DoRequest[[]*Institution](ctx, s.client, http.MethodPost, "merchant/ecobankafrica/institutions", opt)
	if err != nil {
		return nil, resp, err
	}

	return *institutions, resp, nil
}

// RemitteeAccount represents the account details of a remittee.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#68970106-787a-4cfe-917f-91b2e3701bf3
type RemitteeAccount struct {
	AccountStatus string `json:"accountStatus"`
	AccountName   string `json:"accountName"`
	AccountType   string `json:"accountType"`
	BranchCode    string `json:"branchCode"`
	AccountNo     string `json:"accountNo"`
	Currency      string `json:"ccy"`
	AffiliateCode string `json:"affiliateCode"`
}

// GetRemitteeAccountOptions specifies the request parameters to get the account details of a remittee.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#68970106-787a-4cfe-917f-91b2e3701bf3
type GetRemitteeAccountOptions struct {
	RequestId             string `json:"requestId"`
	ClientId              string `json:"clientId"`
	AffiliateCode         string `json:"affiliateCode"`
	DeliveryMethod        string `json:"deliveryMethod"`
	DestinationEntityCode string `json:"destinationEntityCode"`
	AccountNo             string `json:"accountNo"`
	DestinationCountry    string `json:"destinationCountry"`

	secureHashOption
}

// GetAccount returns account details of a supplied account.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#68970106-787a-4cfe-917f-91b2e3701bf3
func (s *RemittanceService) GetAccount(ctx context.Context, opt *GetRemitteeAccountOptions) (*RemitteeAccount, *Response, error) {
	return DoRequest[RemitteeAccount](ctx, s.client, http.MethodPost, "merchant/ecobankafrica/account/enquiry", opt)
}

// Pay is a wrapper around the PaymentService.Pay method.
func (s *RemittanceService) Pay(ctx context.Context, opt *PaymentOptions) (*string, *Response, error) {
	return s.client.Payment.Pay(ctx, opt)
}
