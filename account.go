package ecobank

import (
	"context"

	"github.com/shopspring/decimal"
)

// AccountService represents the account service of the Ecobank API.
type AccountService struct {
	client *Client
}

// AccountBalanceResponse represents a response to an account balance request.
type AccountBalanceResponse struct {
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

// AccountEnquiry represents a response to an account enquiry request.
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

// AccountEnquiryThirdParty represents the response from the account inquiry for third-party payment.
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

// StatementTransaction represents a single transaction record.
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

// AccountBalanceOptions represents a request to get account balance.
type AccountBalanceOptions struct {
	RequestID     string `json:"requestId"`
	AffiliateCode string `json:"affiliateCode"`
	AccountNo     string `json:"accountNo"`
	ClientID      string `json:"clientId"`
	CompanyName   string `json:"companyName"`

	secureHashOption
}

// GetBalance gets the account balance for the given account.
func (a *AccountService) GetBalance(ctx context.Context, opt *AccountBalanceOptions) (*AccountBalanceResponse, *Response, error) {
	return DoRequest[AccountBalanceResponse](ctx, a.client, "POST", "merchant/accountbalance", opt)
}

// AccountEnquiryOptions represents a request to get account details.
type AccountEnquiryOptions struct {
	RequestID     string `json:"requestId"`
	AffiliateCode string `json:"affiliateCode"`
	AccountNo     string `json:"accountNo"`
	ClientID      string `json:"clientId"`
	CompanyName   string `json:"companyName"`

	secureHashOption
}

// Enquiry gets the account details for the given account.
func (a *AccountService) Enquiry(ctx context.Context, opt *AccountEnquiryOptions) (*AccountEnquiry, *Response, error) {
	return DoRequest[AccountEnquiry](ctx, a.client, "POST", "merchant/accountinquiry", opt)
}

// AccountEnquiryThirdPartyOptions represents a request to perform an account inquiry for third-party payment.
type AccountEnquiryThirdPartyOptions struct {
	RequestID           string `json:"requestId"`
	AffiliateCode       string `json:"affiliateCode"`
	AccountNo           string `json:"accountNo"`
	DestinationBankCode string `json:"destinationBankCode"`
	ClientID            string `json:"clientId"`
	CompanyName         string `json:"companyName"`
	SecureHash          string `json:"secureHash"`
}

// EnquiryThirdParty performs an account inquiry for third-party payment.
func (a *AccountService) EnquiryThirdParty(ctx context.Context, opt *AccountEnquiryThirdPartyOptions) (*AccountEnquiryThirdParty, *Response, error) {
	return DoRequest[AccountEnquiryThirdParty](ctx, a.client, "POST", "merchant/accountinquirythridpay", opt)
}

// GenerateStatementOptions represents a request to generate an account statement.
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
func (a *AccountService) GenerateStatement(ctx context.Context, opt *GenerateStatementOptions) ([]*StatementTransaction, *Response, error) {
	statements, resp, err := DoRequest[[]*StatementTransaction](ctx, a.client, "POST", "merchant/statement", opt)
	if err != nil || statements == nil {
		return nil, resp, err
	}

	return *statements, resp, nil
}
