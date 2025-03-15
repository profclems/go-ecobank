package ecobank

import (
	"encoding/json"
	"fmt"
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

type PaymentParamInterface interface {
	json.Marshaler
}

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

// MarshalJSON implements json.Marshaler interface
func (param *PaymentParams[T]) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(param.param)
	if err != nil {
		return nil, err
	}

	// converted to a single quoted string
	return []byte(fmt.Sprintf("%q", b)), nil
}

// DomesticTransferParams represents the parameters for a domestic transfer.
type DomesticTransferParams struct {
	CreditAccountNo     string `json:"creditAccountNo"`
	DebitAccountBranch  string `json:"debitAccountBranch"`
	DebitAccountType    string `json:"debitAccountType"`
	CreditAccountBranch string `json:"creditAccountBranch"`
	CreditAccountType   string `json:"creditAccountType"`
	Amount              string `json:"amount"`
	Currency            string `json:"ccy"`
}

// TokenTransferParams represents the parameters for a token transfer.
type TokenTransferParams struct {
	TransactionDescription string `json:"transactionDescription"`
	SecretCode             string `json:"secretCode"`
	SourceAccount          string `json:"sourceAccount"`
	SourceAccountCurrency  string `json:"sourceAccountCurrency"`
	SourceAccountType      string `json:"sourceAccountType"`
	SenderName             string `json:"senderName"`
	Currency               string `json:"ccy"`
	SenderMobileNo         string `json:"senderMobileNo"`
	Amount                 string `json:"amount"`
	SenderID               string `json:"senderId"`
	BeneficiaryName        string `json:"beneficiaryName"`
	BeneficiaryMobileNo    string `json:"beneficiaryMobileNo"`
	WithdrawalChannel      string `json:"withdrawalChannel"`
}

// TokenIAParams represents the parameters for a token inter-affiliate transfer.
type TokenIAParams struct {
	DestAffiliate          string `json:"destAffiliate"`
	DestCrncy              string `json:"destCrncy"`
	DestinationAccount     string `json:"destinationAccount"`
	DestinationAccountName string `json:"destinationAccountName"`
	ReceiveFirstName       string `json:"receiveFirstName"`
	ReceiveLastName        string `json:"receiveLastName"`
	ReceiverPhoneNumber    string `json:"receiverPhoneNumber"`
	ReceiveEmailAddress    string `json:"receiveEmailAddress"`
	ReceiveIdType          string `json:"receiveIdType"`
	ReceiveIdNumber        string `json:"receiveIdNumber"`
	SourceAmount           string `json:"sourceAmount"`
	TestQuestion           string `json:"testQuestion"`
	TestAnswer             string `json:"testAnswer"`
	Narration              string `json:"narration"`
	PurposeOfTransfer      string `json:"purposeOfTransfer"`
	SendExternalRef        string `json:"sendExternalRef"`
}

// InterbankTransferParams represents the parameters for an interbank transfer.
type InterbankTransferParams struct {
	DestinationBankCode  string `json:"destinationBankCode"`
	SenderName           string `json:"senderName"`
	SenderAddress        string `json:"senderAddress"`
	SenderPhone          string `json:"senderPhone"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo"`
	BeneficiaryName      string `json:"beneficiaryName"`
	BeneficiaryPhone     string `json:"beneficiaryPhone"`
	TransferReferenceNo  string `json:"transferReferenceNo"`
	Amount               string `json:"amount"`
	Currency             string `json:"ccy"`
	TransferType         string `json:"transferType"`
}

// InterbankIAParams represents the parameters for an interbank inter-affiliate transfer.
type InterbankIAParams struct {
	DestinationCountry   string `json:"destinationCountry"`
	DestinationBankCode  string `json:"destinationBankCode"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo"`
	BeneficiaryName      string `json:"beneficiaryName"`
	BeneficiaryPhone     string `json:"beneficiaryPhone"`
	Amount               string `json:"amount"`
	TransferCurrency     string `json:"transferCurrency"`
	TransferReason       string `json:"transferReason"`
	SettleCurrency       string `json:"settleCurrency"`
}

// BillPaymentParams represents the parameters for a bill payment.
type BillPaymentParams struct {
	BillerCode    string             `json:"billerCode"`
	BillRefNo     string             `json:"billRefNo"`
	CbaRefNo      string             `json:"cbaRefNo"`
	CustomerName  string             `json:"customerName"`
	CustomerRefNo string             `json:"customerRefNo"`
	ProductCode   string             `json:"productCode"`
	FormDataValue []BillPaymentField `json:"formDataValue"`
}

// BillPaymentField represents a field in the bill payment form.
type BillPaymentField struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

// AirtimeTopupParams represents the parameters for an airtime topup.
type AirtimeTopupParams struct {
	BillerCode    string              `json:"billerCode"`
	BillRefNo     string              `json:"billRefNo"`
	CbaRefNo      string              `json:"cbaRefNo"`
	CustomerName  string              `json:"customerName"`
	CustomerRefNo string              `json:"customerRefNo"`
	ProductCode   string              `json:"productCode"`
	FormDataValue []AirtimeTopupField `json:"formDataValue"`
}

// AirtimeTopupField represents a field in the airtime topup form.
type AirtimeTopupField struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

// MomoParams represents the parameters for a mobile money transfer.
type MomoParams struct {
	BillerCode    string      `json:"billerCode"`
	BillRefNo     string      `json:"billRefNo"`
	CbaRefNo      string      `json:"cbaRefNo"`
	CustomerName  string      `json:"customerName"`
	CustomerRefNo string      `json:"customerRefNo"`
	ProductCode   string      `json:"productCode"`
	FormDataValue []MomoField `json:"formDataValue"`
}

// MomoField represents a field in the mobile money transfer form.
type MomoField struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}
