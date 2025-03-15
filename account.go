package ecobank

import (
	"context"

	"github.com/shopspring/decimal"
)

// AccountService handles communication with the account related methods of the Ecobank API.
// It combines both account services and account opening services.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#03be64b9-f0dd-4df6-9cca-9e8062943bae
type AccountService struct {
	client *Client
}

// AccountBalance represents a response to an account balance request.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#89d7f8b9-49d8-4a8a-ae3c-acd26cb3e6fe
type AccountBalance struct {
	HostHeaderInfo struct {
		SourceCode      string `json:"sourceCode"`
		RequestID       string `json:"requestId"`
		AffiliateCode   string `json:"affiliateCode"`
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	} `json:"hostHeaderInfo"`

	AccountNo        string          `json:"accountNo"`
	ResponseCode     string          `json:"responseCode"`
	ResponseMessage  string          `json:"responseMessage"`
	AccountName      string          `json:"accountName"`
	Currency         string          `json:"ccy"`
	BranchCode       string          `json:"branchCode"`
	CustomerID       string          `json:"customerID"`
	AvailableBalance decimal.Decimal `json:"availableBalance"`
	CurrentBalance   decimal.Decimal `json:"currentBalance"`
	OverdraftLimit   decimal.Decimal `json:"odlimit"`
	AccountType      string          `json:"accountType"`
	AccountClass     string          `json:"accountClass"`
	AccountStatus    string          `json:"accountStatus"`
}

// AccountBalanceOptions represents a request to get account balance.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#89d7f8b9-49d8-4a8a-ae3c-acd26cb3e6fe
type AccountBalanceOptions struct {
	RequestID     string `json:"requestId"`
	AffiliateCode string `json:"affiliateCode"`
	AccountNo     string `json:"accountNo"`
	ClientID      string `json:"clientId"`
	CompanyName   string `json:"companyName"`

	secureHashOption
}

// GetBalance gets the account balance for the given account.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#89d7f8b9-49d8-4a8a-ae3c-acd26cb3e6fe
func (a *AccountService) GetBalance(ctx context.Context, opt *AccountBalanceOptions) (*AccountBalance, *Response, error) {
	return DoRequest[AccountBalance](ctx, a.client, "POST", "merchant/accountbalance", opt)
}

// AccountEnquiry represents a response to an account enquiry request.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#065afcf7-402b-4625-82d2-24f2dbbfe663
type AccountEnquiry struct {
	AccountNo       string `json:"accountNo"`
	AccountName     string `json:"accountName"`
	Currency        string `json:"ccy"`
	AccountStatus   string `json:"accountStatus"`
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	AffiliateCode   string `json:"affiliateCode"`
	RequestID       string `json:"requestId"`
	SourceCode      string `json:"sourceCode"`
}

// AccountEnquiryOptions represents a request to get account details.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#065afcf7-402b-4625-82d2-24f2dbbfe663
type AccountEnquiryOptions struct {
	RequestID     string `json:"requestId"`
	AffiliateCode string `json:"affiliateCode"`
	AccountNo     string `json:"accountNo"`
	ClientID      string `json:"clientId"`
	CompanyName   string `json:"companyName"`

	secureHashOption
}

// Enquiry gets the account details for the given account.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#065afcf7-402b-4625-82d2-24f2dbbfe663
func (a *AccountService) Enquiry(ctx context.Context, opt *AccountEnquiryOptions) (*AccountEnquiry, *Response, error) {
	return DoRequest[AccountEnquiry](ctx, a.client, "POST", "merchant/accountinquiry", opt)
}

// AccountEnquiryThirdParty represents the response from the account inquiry for third-party payment.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#26923112-e8b8-4956-9f64-0f7f7b489290
type AccountEnquiryThirdParty struct {
	AccountName    string `json:"accountName"`
	AccountType    string `json:"accountType"`
	AccountStatus  string `json:"accountStatus"`
	HostHeaderInfo struct {
		SourceCode      string `json:"sourceCode"`
		RequestID       string `json:"requestId"`
		AffiliateCode   string `json:"affiliateCode"`
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	} `json:"hostHeaderInfo"`
}

// AccountEnquiryThirdPartyOptions represents a request to perform an account inquiry for third-party payment.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#26923112-e8b8-4956-9f64-0f7f7b489290
type AccountEnquiryThirdPartyOptions struct {
	RequestID           string `json:"requestId"`
	AffiliateCode       string `json:"affiliateCode"`
	AccountNo           string `json:"accountNo"`
	DestinationBankCode string `json:"destinationBankCode"`
	ClientID            string `json:"clientId"`
	CompanyName         string `json:"companyName"`

	secureHashOption
}

// EnquiryThirdParty performs an account inquiry for third-party payment.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#26923112-e8b8-4956-9f64-0f7f7b489290
func (a *AccountService) EnquiryThirdParty(ctx context.Context, opt *AccountEnquiryThirdPartyOptions) (*AccountEnquiryThirdParty, *Response, error) {
	return DoRequest[AccountEnquiryThirdParty](ctx, a.client, "POST", "merchant/accountinquirythridpay", opt)
}

// StatementTransaction represents a single transaction record.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#01be6373-e019-4995-aca3-733366acf557
type StatementTransaction struct {
	AccCurrency string `json:"acccy"`
	DebitCredit string `json:"drcrind"`
	RefNumber   string `json:"trnrefno"`
	PaidIn      string `json:"paidin,omitempty"`
	PaidOut     string `json:"paidout,omitempty"`
	ValueDate   Time   `json:"valuedate"`
	Amount      string `json:"lcyamount1"`
	Narrative   string `json:"narrative"`
}

// GenerateStatementOptions represents a request to generate an account statement.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#01be6373-e019-4995-aca3-733366acf557
type GenerateStatementOptions struct {
	CorporateID   string `json:"corporateId"`
	RequestID     string `json:"requestId"`
	ClientID      string `json:"clientId"`
	AffiliateCode string `json:"affiliateCode"`
	AccountNumber string `json:"accountNumber"`
	StartDate     Date   `json:"startDate"`
	EndDate       Date   `json:"endDate"`

	secureHashOption
}

// GenerateStatement generates an account statement for the given account.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#01be6373-e019-4995-aca3-733366acf557
func (a *AccountService) GenerateStatement(ctx context.Context, opt *GenerateStatementOptions) ([]*StatementTransaction, *Response, error) {
	statements, resp, err := DoRequest[[]*StatementTransaction](ctx, a.client, "POST", "merchant/statement", opt)
	if err != nil || statements == nil {
		return nil, resp, err
	}

	return *statements, resp, nil
}

// CreateAccountOptions represents the parameters for creating an account.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#80dc2169-8b2c-435e-8259-5bda0f6ab94c
type CreateAccountOptions struct {
	ClientID           string `json:"clientId"`
	RequestID          string `json:"requestId"`
	AffiliateCode      string `json:"affiliateCode"`
	FirstName          string `json:"firstName"`
	Middlename         string `json:"middlename"`
	Lastname           string `json:"lastname"`
	MobileNo           string `json:"mobileNo"`
	Gender             string `json:"gender"`
	IdentityNo         string `json:"identityNo"`
	IdentityType       string `json:"identityType"`
	IDIssueDate        string `json:"iDIssueDate"`
	IDExpiryDate       string `json:"iDExpiryDate"`
	Ccy                string `json:"ccy"`
	Country            string `json:"country"`
	BranchCode         string `json:"branchCode"`
	DateOfBirth        string `json:"dateOfBirth"`
	CountryOfResidence string `json:"countryOfResidence"`
	Email              string `json:"email"`
	Street             string `json:"street"`
	City               string `json:"city"`
	State              string `json:"state"`
	Image              string `json:"image"`
	Signature          string `json:"signature"`

	secureHashOption
}

// CreateAccountResponse represents the response after attempting to create an account.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#80dc2169-8b2c-435e-8259-5bda0f6ab94c
type CreateAccountResponse struct {
	Shortname      string `json:"shortname"`
	AccountNo      string `json:"accountNo"`
	MobileNo       string `json:"mobileNo"`
	TrackRef       string `json:"trackRef"`
	ClientID       string `json:"clientId"`
	HostHeaderInfo struct {
		SourceCode      string `json:"sourceCode"`
		RequestID       string `json:"requestId"`
		AffiliateCode   string `json:"affiliateCode"`
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	} `json:"hostHeaderInfo"`
}

// CreateAccount creates an account.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#80dc2169-8b2c-435e-8259-5bda0f6ab94c
func (a *AccountService) CreateAccount(ctx context.Context, opt *CreateAccountOptions) (*CreateAccountResponse, *Response, error) {
	return DoRequest[CreateAccountResponse](ctx, a.client, "POST", "merchant/createexpressaccount", opt)
}
