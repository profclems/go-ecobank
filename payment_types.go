package ecobank

import "encoding/json"

type PaymentType string

const (
	BILLPAYMENT  PaymentType = "BILLPAYMENT"
	TOKEN        PaymentType = "TOKEN"
	DOMESTIC     PaymentType = "DOMESTIC"
	INTERBANK    PaymentType = "INTERBANK"
	INTEBBANKIA  PaymentType = "INTEBBANKIA"
	AIRTIMETOPUP PaymentType = "AIRTIMETOPUP"
	MOMO         PaymentType = "MOMO"
)

type SupportedPaymentParamType interface {
	json.Unmarshaler
}

type BillPaymentParams struct {
}
