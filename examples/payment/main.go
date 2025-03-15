package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/profclems/go-ecobank"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()

	username := os.Getenv("ECOBANK_USERNAME")
	password := os.Getenv("ECOBANK_PASSWORD")
	labKey := os.Getenv("ECOBANK_LAB_KEY")

	client, err := ecobank.NewClient(username, password, labKey)
	checkErr(errors.Wrap(err, "failed to initiate client"))

	err = client.Login(ctx)
	checkErr(errors.Wrap(err, "failed to login"))

	req := &ecobank.PaymentOptions{
		PaymentHeader: ecobank.PaymentHeader{
			Clientid:          "EGHTelc000043",
			Batchsequence:     "1",
			Batchamount:       decimal.NewFromInt(520),
			Transactionamount: decimal.NewFromInt(520),
			Batchid:           "EG1593490",
			Transactioncount:  6,
			Batchcount:        6,
			Transactionid:     "E12T443308",
			Debittype:         "Multiple",
			AffiliateCode:     "EGH",
			Totalbatches:      "1",
			ExecutionDate:     ecobank.NewTime(time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)),
		},
		Extension: []ecobank.PaymentExtension{
			{
				RequestId:   "2323",
				RequestType: ecobank.DOMESTIC,
				ParamList: ecobank.NewPaymentParams(ecobank.DomesticTransferParams{
					CreditAccountNo:     "1441001996321",
					DebitAccountBranch:  "ACCRA",
					DebitAccountType:    "Corporate",
					CreditAccountBranch: "Accra",
					CreditAccountType:   "Corporate",
					Amount:              "10",
					Currency:            "GHS",
				}),
				Amount:   decimal.NewFromInt(10),
				Currency: "GHS",
				RateType: "spot",
			},
			{
				RequestId:   "432",
				RequestType: ecobank.TOKEN,
				ParamList: ecobank.NewPaymentParams(ecobank.TokenTransferParams{
					TransactionDescription: "Service payment for electrical repairs.",
					SecretCode:             "AWER1234",
					SourceAccount:          "1441000565307",
					SourceAccountCurrency:  "GHS",
					SourceAccountType:      "Corporate",
					SenderName:             "Freeman Kay",
					Currency:               "GHS",
					SenderMobileNo:         "0202205113",
					Amount:                 "40",
					SenderID:               "QWE345Y4",
					BeneficiaryName:        "Stephen Kojo",
					BeneficiaryMobileNo:    "0233445566",
					WithdrawalChannel:      "ATM",
				}),
				Amount:   decimal.NewFromInt(40),
				Currency: "GHS",
				RateType: "spot",
			},
			{
				RequestId:   "2325",
				RequestType: ecobank.INTERBANK,
				ParamList: ecobank.NewPaymentParams(ecobank.InterbankTransferParams{
					DestinationBankCode:  "ASB",
					SenderName:           "BEN",
					SenderAddress:        "23 Accra Central",
					SenderPhone:          "233263653712",
					BeneficiaryAccountNo: "110424812001",
					BeneficiaryName:      "Owen",
					BeneficiaryPhone:     "233543837123",
					TransferReferenceNo:  "QWE345Y4",
					Amount:               "10",
					Currency:             "GHS",
					TransferType:         "spot",
				}),
				Amount:   decimal.NewFromInt(10),
				Currency: "GHS",
				RateType: "spot",
			},
			{
				RequestId:   "ECI55096987905",
				RequestType: ecobank.BILLPAYMENT,
				ParamList: ecobank.NewPaymentParams(ecobank.BillPaymentParams{
					BillerCode:    "Pass_Bio_ECI",
					BillRefNo:     "239729",
					CbaRefNo:      "",
					CustomerName:  "Freeman Kay",
					CustomerRefNo: "239729",
					ProductCode:   "PassBio",
					FormDataValue: []ecobank.BillPaymentField{
						{FieldName: "LastName", FieldValue: "Kojo"},
						{FieldName: "FirstName", FieldValue: "Kwame"},
						{FieldName: "Amount", FieldValue: "300"},
						{FieldName: "Phone", FieldValue: "225543756765"},
						{FieldName: "Email", FieldValue: "enyaledzigbor@ecobank.com"},
						{FieldName: "reference", FieldValue: "210120400582"},
					},
				}),
				Amount:   decimal.NewFromInt(300),
				Currency: "GHS",
				RateType: "spot",
			},
			{
				RequestId:   "WQ5500098663046",
				RequestType: ecobank.AIRTIMETOPUP,
				ParamList: ecobank.NewPaymentParams(ecobank.AirtimeTopupParams{
					BillerCode:    "A02E",
					BillRefNo:     "81729",
					CbaRefNo:      "",
					CustomerName:  "Owen Kay",
					CustomerRefNo: "824225",
					ProductCode:   "A02E",
					FormDataValue: []ecobank.AirtimeTopupField{
						{FieldName: "BEN_PHONE_NO", FieldValue: "2348034830707"},
					},
				}),
				Amount:   decimal.NewFromInt(10),
				Currency: "NGN",
				RateType: "spot",
			},
			{
				RequestId:   "1234BBY8SXZX",
				RequestType: ecobank.MOMO,
				ParamList: ecobank.NewPaymentParams(ecobank.MomoParams{
					BillerCode:    "AIRTELTIGOEGH",
					BillRefNo:     "2988759",
					CbaRefNo:      "05609",
					CustomerName:  "Owen Kay",
					CustomerRefNo: "824225",
					ProductCode:   "AIRTELTIGO_MOBILEMONEY",
					FormDataValue: []ecobank.MomoField{
						{FieldName: "BEN_PHONE_NO", FieldValue: "0560000159"},
					},
				}),
				Amount:   decimal.NewFromInt(150),
				Currency: "GHS",
				RateType: "spot",
			},
			{
				RequestId:   "ByA1lUeHKK",
				RequestType: ecobank.INTERBANKIA,
				ParamList: ecobank.NewPaymentParams(ecobank.InterbankIAParams{
					DestinationCountry:   "EKE",
					DestinationBankCode:  "11000",
					BeneficiaryAccountNo: "01108738185300",
					BeneficiaryName:      "Harmony Blinks",
					BeneficiaryPhone:     "233543837123",
					Amount:               "100",
					TransferCurrency:     "GHS",
					TransferReason:       "vacation",
					SettleCurrency:       "USD",
				}),
				Amount:   decimal.NewFromInt(100),
				Currency: "GHS",
				RateType: "spot",
			},
			{
				RequestId:   "ByA1lUeHUU",
				RequestType: ecobank.TOKENIA,
				ParamList: ecobank.NewPaymentParams(ecobank.TokenIAParams{
					DestAffiliate:          "EKE",
					DestCrncy:              "KES",
					DestinationAccount:     "01108738185300",
					DestinationAccountName: "Owen",
					ReceiveFirstName:       "Eugene",
					ReceiveLastName:        "Eugene",
					ReceiverPhoneNumber:    "233543837123",
					ReceiveEmailAddress:    "test@gmail.com",
					ReceiveIdType:          "23421221",
					ReceiveIdNumber:        "23421221",
					SourceAmount:           "100",
					TestQuestion:           "GHS",
					TestAnswer:             "vacation",
					Narration:              "narration",
					PurposeOfTransfer:      "purposeOfTransfer",
					SendExternalRef:        "ByA1lUeHUU",
				}),
				Amount:   decimal.NewFromInt(100),
				Currency: "GHS",
				RateType: "spot",
			},
			{
				RequestId:   "ByA1lUeHSS",
				RequestType: ecobank.TOKENIA,
				ParamList: ecobank.NewPaymentParams(ecobank.TokenIAParams{
					DestAffiliate:          "EKE",
					DestCrncy:              "KES",
					DestinationAccount:     "01108738185300",
					DestinationAccountName: "Owen",
					ReceiveFirstName:       "Eugene",
					ReceiveLastName:        "Eugene",
					ReceiverPhoneNumber:    "233543837123",
					ReceiveEmailAddress:    "test@gmail.com",
					ReceiveIdType:          "23421221",
					ReceiveIdNumber:        "23421221",
					SourceAmount:           "100",
					TestQuestion:           "GHS",
					TestAnswer:             "vacation",
					Narration:              "narration",
					PurposeOfTransfer:      "purposeOfTransfer",
					SendExternalRef:        "ByA1lUeHSS",
				}),
				Amount:   decimal.NewFromInt(100),
				Currency: "GHS",
				RateType: "spot",
			},
		},
	}

	req.SetHash("398d4f285cc33e12f035da19fa9d954be35afaf66816531c4f1a1aedd3c6f132a85c62b23ca12d7b9a99bf5a84fc69b66738289a70e8f8115e90ffaa060f4026")

	// make payment
	paymentStatus, resp, err := client.Payment.Pay(ctx, req)

	checkErr(errors.Wrap(err, "failed to make payment"))

	fmt.Printf("Payment status: %v\n", *paymentStatus)
	fmt.Printf("Code: %+v\n", resp.Code)
	fmt.Printf("Message: %+v\n", resp.Message)
	fmt.Printf("HTTP Status: %+v\n", resp.Status)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
