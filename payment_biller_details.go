package ecobank

import "github.com/shopspring/decimal"

// BillFormData represents input fields required for billing.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#22c57a29-be69-4ca6-8274-896defa6b2f9
type BillFormData struct {
	SerialNo       int      `json:"serialNo"`
	FieldName      string   `json:"fieldName"`
	FieldTitle     string   `json:"fieldTitle"`
	DataType       string   `json:"dataType"`
	ValidateField  string   `json:"validateField"`
	DefaultValue   string   `json:"defaultValue"`
	MaxFieldLength int      `json:"maxFieldLength"`
	ListOfValues   string   `json:"listofValues"`
	LookupValue    []string `json:"lookupValue"`
}

// BillerProductInfo represents biller product details.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#22c57a29-be69-4ca6-8274-896defa6b2f9
type BillerProductInfo struct {
	ProductCode        string          `json:"productCode"`
	ProductName        string          `json:"productName"`
	ProductDescription string          `json:"productDescription"`
	ProductCategory    string          `json:"productCategory"`
	AmountType         string          `json:"amountType"`
	MinAmount          decimal.Decimal `json:"minAmount"`
	MaxAmount          decimal.Decimal `json:"maxAmount"`
	Currency           string          `json:"ccy"`
	ExchangeRate       decimal.Decimal `json:"exchRate"`
}

// BillerDetails represents the response payload for retrieving biller details.
//
// API docs: https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#22c57a29-be69-4ca6-8274-896defa6b2f9
type BillerDetails struct {
	BillerInfo struct {
		BillerCode                string          `json:"billerCode"`
		BillerID                  int             `json:"billerID"`
		BillerName                string          `json:"billerName"`
		BillerDescription         string          `json:"billerDescription"`
		BillerCategory            string          `json:"billerCategory"`
		BillerEmail               string          `json:"billerEmail"`
		BillerPhone               string          `json:"billerPhone"`
		BillerSite                string          `json:"billerSite"`
		BillerLogo                string          `json:"billerLogo"`
		BillAmountType            string          `json:"billAmountType"`
		BillAmount                decimal.Decimal `json:"billAmount"`
		CollectionAccountNo       string          `json:"collectionAccountNo"`
		CollectionAccountName     string          `json:"collectionAccountName"`
		CollectionAccountBankCode string          `json:"collectionAccountBankCode"`
		AggregatorName            string          `json:"aggregatorName"`
		ValidationRequired        string          `json:"validationRequired"`
		ProductList               string          `json:"productList"`
	} `json:"billerDetail"`

	BillFormData []BillFormData `json:"billFormData"`

	BillerProductInfo []BillerProductInfo `json:"billerProductInfo"`
	HostHeaderInfo    struct {
		SourceCode      string `json:"sourceCode"`
		RequestID       string `json:"requestId"`
		AffiliateCode   string `json:"affiliateCode"`
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	} `json:"hostHeaderInfo"`
}
