package ecobank

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
)

// PaymentType represents the supported payment types for the Ecobank API.
//
// API docs:
//   - https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#fca97841-db96-4828-bc1b-525e973efe91
//   - https://documenter.getpostman.com/view/9576712/2s7YtWCtNX#acfe7f55-27aa-487d-ba1a-799ecb466bd7
type PaymentType string

const (
	// BILLPAYMENT allows customer to make bill payment, by enabling customer
	// to validate references pertaining to a bill and further allow payment.
	BILLPAYMENT PaymentType = "BILLPAYMENT"

	// TOKEN allows customer to make cash payment to their beneficiaries and
	// enable withdrawal to be done via ATM or Xpress point using the
	// generated Token upon successfully sending an instruction.
	TOKEN PaymentType = "TOKEN"

	// TOKENIA is used for movement of funds between countries to credit
	// a third party or local bank account at the receiving country/affiliate.
	// For example, sending from an Ecobank Ghana account to an account
	// in another bank (other than Ecobank) in the receiving country.
	TOKENIA PaymentType = "TOKENIA"

	// DOMESTIC enables customer to send an instruction to debit an account
	// within Ecobank and credit another Ecobank account holder of the same country
	DOMESTIC PaymentType = "DOMESTIC"

	// INTERBANK enables customer debit an account within Ecobank and credit
	// a beneficiary's account with another bank in the same country.
	// This caters for Instant payment, RTGS, SICA, SYGMA and TFT transactions.
	INTERBANK PaymentType = "INTERBANK"

	// INTERBANKIA is used for movement of funds between countries to credit
	// a third party or local bank account at the receiving country/affiliate.
	// For example, sending from an Ecobank Ghana account to an account
	// in another bank (other than Ecobank) in the receiving country.
	INTERBANKIA PaymentType = "INTERBANKIA"

	// AIRTIMETOPUP allows customer make airtime top-up request to credit Momo accounts
	// across affiliates where Ecobank have an integration with MNO’s
	AIRTIMETOPUP PaymentType = "AIRTIMETOPUP"

	// MOMO allows customers to initiate bank-to-wallet transactions across
	// affiliates where Ecobank has an integration with the telcos/aggregators
	MOMO PaymentType = "MOMO"
	// MOMOIA is used for the movement of funds between countries to credit
	// a mobile wallet account in the receiving affiliate/country.
	MOMOIA PaymentType = "MOMOIA"
)

// PaymentParamInterface defines an interface for payment parameters that can be serialized into JSON.
type PaymentParamInterface interface {
	json.Marshaler
}

// SupportedPaymentParamTypes is a type constraint for generics in PaymentParams.
// It ensures that only predefined payment parameter structs can be used.
type SupportedPaymentParamTypes interface {
	DomesticTransferParams | TokenTransferParams | InterbankTransferParams |
		BillPaymentParams | AirtimeTopupParams | MomoParams |
		TokenIAParams | InterbankIAParams
}

// PaymentParams represents the parameters for a payment.
type PaymentParams[T SupportedPaymentParamTypes] struct {
	param T
}

// NewPaymentParams creates a new PaymentParams.
func NewPaymentParams[T SupportedPaymentParamTypes](param T) PaymentParamInterface {
	return &PaymentParams[T]{param: param}
}

// MarshalJSON implements the json.Marshaler interface for PaymentParams.
// It serializes the struct fields into a JSON-encoded string that follows
// a specific key-value format required by the ecobank API.
//
// **Format:**
// The output is a JSON array of objects, where each object represents a field
// with two properties:
//   - "key": The JSON field name of the struct.
//   - "value": The corresponding value from the struct.
//
// Example Output:
// ```json
// "[
//
//	{\"key\": \"billerCode\", \"value\": \"Pass_Bio_ECI\"},
//	{\"key\": \"billRefNo\", \"value\": \"239729\"},
//	{\"key\": \"customerName\", \"value\": \"Freeman Kay\"},
//	{\"key\": \"formDataValue\", \"value\": \"[{\\\"fieldName\\\": \\\"LastName\\\", \\\"fieldValue\\\": \\\"Kojo\\\"}]\"}
//
// ]"
// ```
//
// **Special Handling:**
//   - The `formDataValue` field is itself a JSON-encoded array of objects, each
//     containing a `fieldName` and `fieldValue`. This nested structure is stringified
//     to maintain consistency with the external API's expected format.
func (param *PaymentParams[T]) MarshalJSON() ([]byte, error) {
	val := reflect.ValueOf(param.param)
	typ := reflect.TypeOf(param.param)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	var b strings.Builder
	b.WriteString(`"[`)

	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		// get the json tag of the field
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// now we need to set the {"key": "jsonTag", "value": "value"}
		b.WriteString(`{\"key\": \"`)
		b.WriteString(jsonTag)
		b.WriteString(`\", \"value\": \"`)
		// check if the field is FormDataArray
		if jsonTag == "formDataValue" {
			formData := fieldValue.Interface().(FormDataArray)
			b.WriteString(`[`)
			for i, fd := range formData {
				b.WriteString(`{\\\"fieldName\\\": \\\"`)
				b.WriteString(fd.FieldName)
				b.WriteString(`\\\", \\\"fieldValue\\\": \\\"`)
				b.WriteString(fd.FieldValue)
				b.WriteString(`\\\"}`)
				if i < len(formData)-1 {
					b.WriteString(`,`)
				}
			}
			b.WriteString(`]`)
		} else {
			b.WriteString(formatToStr(fieldValue.Interface()))

		}
		b.WriteString(`\"}`)
		if i < val.NumField()-1 {
			b.WriteString(`,`)
		}
	}

	b.WriteString(`]"`)

	// converted to a single quoted string
	return []byte(b.String()), nil
}

// DomesticTransferParams represents the parameters for DOMESTIC payment type.
type DomesticTransferParams struct {
	CreditAccountNo     string          `json:"creditAccountNo"`
	DebitAccountBranch  string          `json:"debitAccountBranch"`
	DebitAccountType    string          `json:"debitAccountType"`
	CreditAccountBranch string          `json:"creditAccountBranch"`
	CreditAccountType   string          `json:"creditAccountType"`
	Amount              decimal.Decimal `json:"amount"`
	Currency            string          `json:"ccy"`
}

// TokenTransferParams represents the parameters TOKEN payment type.
type TokenTransferParams struct {
	TransactionDescription string          `json:"transactionDescription"`
	SecretCode             string          `json:"secretCode"`
	SourceAccount          string          `json:"sourceAccount"`
	SourceAccountCurrency  string          `json:"sourceAccountCurrency"`
	SourceAccountType      string          `json:"sourceAccountType"`
	SenderName             string          `json:"senderName"`
	Currency               string          `json:"ccy"`
	SenderMobileNo         string          `json:"senderMobileNo"`
	Amount                 decimal.Decimal `json:"amount"`
	SenderID               string          `json:"senderId"`
	BeneficiaryName        string          `json:"beneficiaryName"`
	BeneficiaryMobileNo    string          `json:"beneficiaryMobileNo"`
	WithdrawalChannel      string          `json:"withdrawalChannel"`
}

// TokenIAParams represents the parameters for TOKENIA payment type.
type TokenIAParams struct {
	DestinationAffiliate   string          `json:"destAffiliate"`
	DestinationCurrency    string          `json:"destCrncy"`
	DestinationAccount     string          `json:"destinationAccount"`
	DestinationAccountName string          `json:"destinationAccountName"`
	ReceiverFirstName      string          `json:"receiveFirstName"`
	ReceiverLastName       string          `json:"receiveLastName"`
	ReceiverPhoneNumber    string          `json:"receiverPhoneNumber"`
	ReceiverEmailAddress   string          `json:"receiveEmailAddress"`
	ReceiverIDType         string          `json:"receiveIdType"`
	ReceiverIDNumber       string          `json:"receiveIdNumber"`
	SourceAmount           decimal.Decimal `json:"sourceAmount"`
	TestQuestion           string          `json:"testQuestion"`
	TestAnswer             string          `json:"testAnswer"`
	Narration              string          `json:"narration"`
	PurposeOfTransfer      string          `json:"purposeOfTransfer"`
	SendExternalRef        string          `json:"sendExternalRef"`
}

// InterbankTransferParams represents the parameters for INTERBANK payment type.
type InterbankTransferParams struct {
	DestinationBankCode  string          `json:"destinationBankCode"`
	SenderName           string          `json:"senderName"`
	SenderAddress        string          `json:"senderAddress"`
	SenderPhone          string          `json:"senderPhone"`
	BeneficiaryAccountNo string          `json:"beneficiaryAccountNo"`
	BeneficiaryName      string          `json:"beneficiaryName"`
	BeneficiaryPhone     string          `json:"beneficiaryPhone"`
	TransferReferenceNo  string          `json:"transferReferenceNo"`
	Amount               decimal.Decimal `json:"amount"`
	Currency             string          `json:"ccy"`
	TransferType         string          `json:"transferType"`
}

// InterbankIAParams represents the parameters for INTERBANKIA payment type.
type InterbankIAParams struct {
	DestinationCountry   string          `json:"destinationCountry"`
	DestinationBankCode  string          `json:"destinationBankCode"`
	BeneficiaryAccountNo string          `json:"beneficiaryAccountNo"`
	BeneficiaryName      string          `json:"beneficiaryName"`
	BeneficiaryPhone     string          `json:"beneficiaryPhone"`
	Amount               decimal.Decimal `json:"amount"`
	TransferCurrency     string          `json:"transferCurrency"`
	TransferReason       string          `json:"transferReason"`
	SettleCurrency       string          `json:"settleCurrency"`
}

// BillPaymentParams represents the parameters for BILLPAYMENT payment type.
type BillPaymentParams struct {
	BillerCode    string        `json:"billerCode"`
	BillRefNo     string        `json:"billRefNo"`
	CbaRefNo      string        `json:"cbaRefNo"`
	CustomerName  string        `json:"customerName"`
	CustomerRefNo string        `json:"customerRefNo"`
	ProductCode   string        `json:"productCode"`
	FormDataValue FormDataArray `json:"formDataValue"`
}

// AirtimeTopupParams represents the parameters for AIRTIMETOPUP payment type.
type AirtimeTopupParams struct {
	BillerCode    string        `json:"billerCode"`
	BillRefNo     string        `json:"billRefNo"`
	CbaRefNo      string        `json:"cbaRefNo"`
	CustomerName  string        `json:"customerName"`
	CustomerRefNo string        `json:"customerRefNo"`
	ProductCode   string        `json:"productCode"`
	FormDataValue FormDataArray `json:"formDataValue"`
}

// MomoParams represents the parameters for MOMO payment type.
type MomoParams struct {
	BillerCode    string        `json:"billerCode"`
	BillRefNo     string        `json:"billRefNo"`
	CbaRefNo      string        `json:"cbaRefNo"`
	CustomerName  string        `json:"customerName"`
	CustomerRefNo string        `json:"customerRefNo"`
	ProductCode   string        `json:"productCode"`
	FormDataValue FormDataArray `json:"formDataValue"`
}

// MomoIAParams represents the parameters for MOMOIA payment type.
type MomoIAParams struct {
	DestinationAffiliate   string          `json:"destAffiliate"`
	DestinationCurrency    string          `json:"destCrncy"`
	DestinationAccount     string          `json:"destinationAccount"`
	DestinationAccountName string          `json:"destinationAccountName"`
	ReceiverFirstName      string          `json:"receiveFirstName"`
	ReceiverLastName       string          `json:"receiveLastName"`
	ReceiverPhoneNumber    string          `json:"receiverPhoneNumber"`
	ReceiverEmailAddress   string          `json:"receiveEmailAddress"`
	ReceiverIDType         string          `json:"receiveIdType"`
	ReceiverIDNumber       string          `json:"receiveIdNumber"`
	SourceAmount           decimal.Decimal `json:"sourceAmount"`
	TestQuestion           string          `json:"testQuestion"`
	TestAnswer             string          `json:"testAnswer"`
	Narration              string          `json:"narration"`
	PurposeOfTransfer      string          `json:"purposeOfTransfer"`
	SendExternalRef        string          `json:"sendExternalRef"`
}

type FormDataArray []FormData

// FormData represents a fieldName and fieldValue pair.
type FormData struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}
