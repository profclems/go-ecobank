package ecobank

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
)

type PaymentType string

const (
	BILLPAYMENT  PaymentType = "BILLPAYMENT"
	TOKEN        PaymentType = "TOKEN"
	TOKENIA      PaymentType = "TOKENIA"
	DOMESTIC     PaymentType = "DOMESTIC"
	INTERBANK    PaymentType = "INTERBANK"
	INTERBANKIA  PaymentType = "INTERBANKIA"
	AIRTIMETOPUP PaymentType = "AIRTIMETOPUP"
	MOMO         PaymentType = "MOMO"
	MOMOIA       PaymentType = "MOMOIA"
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

// DomesticTransferParams represents the parameters for a domestic transfer.
type DomesticTransferParams struct {
	CreditAccountNo     string          `json:"creditAccountNo"`
	DebitAccountBranch  string          `json:"debitAccountBranch"`
	DebitAccountType    string          `json:"debitAccountType"`
	CreditAccountBranch string          `json:"creditAccountBranch"`
	CreditAccountType   string          `json:"creditAccountType"`
	Amount              decimal.Decimal `json:"amount"`
	Currency            string          `json:"ccy"`
}

// TokenTransferParams represents the parameters for a token transfer.
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

// TokenIAParams represents the parameters for a token inter-affiliate transfer.
type TokenIAParams struct {
	DestAffiliate          string          `json:"destAffiliate"`
	DestCrncy              string          `json:"destCrncy"`
	DestinationAccount     string          `json:"destinationAccount"`
	DestinationAccountName string          `json:"destinationAccountName"`
	ReceiveFirstName       string          `json:"receiveFirstName"`
	ReceiveLastName        string          `json:"receiveLastName"`
	ReceiverPhoneNumber    string          `json:"receiverPhoneNumber"`
	ReceiveEmailAddress    string          `json:"receiveEmailAddress"`
	ReceiveIdType          string          `json:"receiveIdType"`
	ReceiveIdNumber        string          `json:"receiveIdNumber"`
	SourceAmount           decimal.Decimal `json:"sourceAmount"`
	TestQuestion           string          `json:"testQuestion"`
	TestAnswer             string          `json:"testAnswer"`
	Narration              string          `json:"narration"`
	PurposeOfTransfer      string          `json:"purposeOfTransfer"`
	SendExternalRef        string          `json:"sendExternalRef"`
}

// InterbankTransferParams represents the parameters for an interbank transfer.
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

// InterbankIAParams represents the parameters for an interbank inter-affiliate transfer.
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

// BillPaymentParams represents the parameters for a bill payment.
type BillPaymentParams struct {
	BillerCode    string        `json:"billerCode"`
	BillRefNo     string        `json:"billRefNo"`
	CbaRefNo      string        `json:"cbaRefNo"`
	CustomerName  string        `json:"customerName"`
	CustomerRefNo string        `json:"customerRefNo"`
	ProductCode   string        `json:"productCode"`
	FormDataValue FormDataArray `json:"formDataValue"`
}

// AirtimeTopupParams represents the parameters for an airtime topup.
type AirtimeTopupParams struct {
	BillerCode    string        `json:"billerCode"`
	BillRefNo     string        `json:"billRefNo"`
	CbaRefNo      string        `json:"cbaRefNo"`
	CustomerName  string        `json:"customerName"`
	CustomerRefNo string        `json:"customerRefNo"`
	ProductCode   string        `json:"productCode"`
	FormDataValue FormDataArray `json:"formDataValue"`
}

// MomoParams represents the parameters for a mobile money transfer.
type MomoParams struct {
	BillerCode    string        `json:"billerCode"`
	BillRefNo     string        `json:"billRefNo"`
	CbaRefNo      string        `json:"cbaRefNo"`
	CustomerName  string        `json:"customerName"`
	CustomerRefNo string        `json:"customerRefNo"`
	ProductCode   string        `json:"productCode"`
	FormDataValue FormDataArray `json:"formDataValue"`
}

type FormDataArray []FormData

// FormData represents a fieldName and fieldValue pair.
type FormData struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}
