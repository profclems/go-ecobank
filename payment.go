package ecobank

import (
	"context"
	"net/http"

	"github.com/shopspring/decimal"
)

// PaymentService represents the payment service of the Ecobank API.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#149a5d48-68d6-459b-92e1-5100607d1311
type PaymentService struct {
	client *Client
}

// BillerInfo represents a single biller entity.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#eec6e30d-de2b-4565-89a1-cded3a7a8284
type BillerInfo struct {
	BillerCode          string          `json:"billerCode"`
	BillerID            int             `json:"billerID"`
	BillerName          string          `json:"billerName"`
	BillerDescription   string          `json:"billerDescription"`
	BillerCategory      string          `json:"billerCategory"`
	BillerLogo          string          `json:"billerLogo"`
	BillAmountType      string          `json:"billAmountType"`
	BillAmount          decimal.Decimal `json:"billAmount"`
	Currency            string          `json:"ccy"`
	CollectionAccountNo string          `json:"collectionAccountNo"`
	AggregatorName      string          `json:"aggregatorName"`
	AmountDenominations string          `json:"amountDenominations"`
	ProductCodeList     string          `json:"productCodeList"`
}

// BillerList is the response payload for getting the biller list.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#eec6e30d-de2b-4565-89a1-cded3a7a8284
type BillerList struct {
	BillerInfo     []BillerInfo `json:"billerInfo"`
	HostHeaderInfo struct {
		SourceCode      string `json:"sourceCode"`
		RequestID       string `json:"requestId"`
		AffiliateCode   string `json:"affiliateCode"`
		ResponseCode    int    `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	} `json:"hostHeaderInfo"`
}

// ValidateBillerResponse represents the response payload for validating a biller.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#575a20cc-d7d1-4627-9665-1211622e1523
type ValidateBillerResponse struct {
	HostHeaderInfo struct {
		SourceCode      string `json:"sourceCode"`
		RequestID       string `json:"requestId"`
		AffiliateCode   string `json:"affiliateCode"`
		ResponseCode    int    `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	} `json:"hostHeaderInfo"`
	BillerCode         string          `json:"billerCode"`
	BillRefNo          string          `json:"billRefNo"`
	CustomerName       string          `json:"customerName"`
	Amount             decimal.Decimal `json:"amount"`
	PaymentDescription string          `json:"paymentDescription"`
	ProductCode        string          `json:"productCode"`
	ResponseValues     string          `json:"responseValues"`
	FormDataValue      []struct {
		FieldName        string `json:"fieldName"`
		FieldDescription string `json:"fieldDescription"`
		FieldMasked      string `json:"fieldMasked"`
		FieldValue       string `json:"fieldValue"`
		FieldRequired    string `json:"fieldRequired"`
		DataType         string `json:"dataType"`
	} `json:"formDataValue"`
}

// GetBillerListOptions represents the request payload for getting the biller list.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#eec6e30d-de2b-4565-89a1-cded3a7a8284
type GetBillerListOptions struct {
	// RequestID identifies the corporation ID provisioned for the corporate
	RequestID string `json:"requestId"`
	// AffiliateCode of which the account and client has been maintained
	AffiliateCode string `json:"affiliateCode"`

	secureHashOption
}

// GetBillerList fetches the list of billers from the Ecobank API.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#eec6e30d-de2b-4565-89a1-cded3a7a8284
func (s *PaymentService) GetBillerList(ctx context.Context, req *GetBillerListOptions) (*BillerList, *Response, error) {
	return DoRequest[BillerList](ctx, s.client, http.MethodPost, "payment/getbillerlist", req)
}

// GetBillerDetailsOptions represents the request payload for retrieving biller details.
type GetBillerDetailsOptions struct {
	RequestID     string `json:"requestId"`
	AffiliateCode string `json:"affiliateCode"`
	BillerCode    string `json:"billerCode"`

	secureHashOption
}

// GetBillerDetails fetches details of a specific biller.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#22c57a29-be69-4ca6-8274-896defa6b2f9
func (p *PaymentService) GetBillerDetails(ctx context.Context, opt *GetBillerDetailsOptions) (*BillerDetails, *Response, error) {
	return DoRequest[BillerDetails](ctx, p.client, http.MethodPost, "/merchant/getbillerdetails", opt)
}

// ValidateBillerOptions represents the request payload for validating a biller.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#575a20cc-d7d1-4627-9665-1211622e1523
type ValidateBillerOptions struct {
	RequestID     string `json:"requestId"`
	AffiliateCode string `json:"affiliateCode"`
	BillerCode    string `json:"billerCode"`
	ProductCode   string `json:"productCode"`
	MobileNumber  string `json:"mobileNnumber"`
	CustomerName  string `json:"customerName"`
	FormDataValue []struct {
		FieldName  string `json:"fieldName"`
		FieldValue string `json:"fieldValue"`
	} `json:"formDataValue"`

	secureHashOption
}

// ValidateBiller validates a biller.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#575a20cc-d7d1-4627-9665-1211622e1523
func (p *PaymentService) ValidateBiller(ctx context.Context, opt *ValidateBillerOptions) (*ValidateBillerResponse, *Response, error) {
	return DoRequest[ValidateBillerResponse](ctx, p.client, http.MethodPost, "/merchant/validatebiller", opt)
}

// PaymentOptions represents a request to make a payment.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#fca97841-db96-4828-bc1b-525e973efe91
type PaymentOptions struct {
	PaymentHeader PaymentHeader      `json:"paymentHeader"`
	Extension     []PaymentExtension `json:"extension"`

	secureHashOption
}

// PaymentHeader is the main payload for a payment request.
type PaymentHeader struct {
	Batchsequence     string          `json:"batchsequence"`
	Batchamount       decimal.Decimal `json:"batchamount"`
	Transactionamount decimal.Decimal `json:"transactionamount"`
	Batchid           string          `json:"batchid"`
	Transactioncount  int             `json:"transactioncount"`
	Batchcount        int             `json:"batchcount"`
	Transactionid     string          `json:"transactionid"`
	Debittype         string          `json:"debittype"`
	AffiliateCode     string          `json:"affiliateCode"`
	Totalbatches      string          `json:"totalbatches"`
	ExecutionDate     Time            `json:"execution_date"`
	Clientid          string          `json:"clientid"`
}

// PaymentExtension contains additional information for a payment request.
type PaymentExtension struct {
	RequestId   string                `json:"request_id"`
	RequestType PaymentType           `json:"request_type"`
	ParamList   PaymentParamInterface `json:"param_list"`
	Amount      decimal.Decimal       `json:"amount"`
	Currency    string                `json:"currency"`
	Status      string                `json:"status"`
	RateType    string                `json:"rate_type"`
}

// Pay sends a payment request to the Ecobank API.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#fca97841-db96-4828-bc1b-525e973efe91
func (p *PaymentService) Pay(ctx context.Context, opt *PaymentOptions) (*string, *Response, error) {
	return DoRequest[string](ctx, p.client, http.MethodPost, "merchant/payment", opt)
}
