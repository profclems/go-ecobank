package ecobank

import (
	"net/http"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccountService_GetBalance(t *testing.T) {
	mockResponse := `{
		"response_code": 200,
		"response_message": "success",
		"response_content": {
			"hostHeaderInfo": {
				"sourceCode": "CORPORATEAPI",
				"requestId": "14232436312",
				"affiliateCode": "EGH",
				"responseCode": "000",
				"responseMessage": "SUCCESS"
			},
			"accountNo": "1441000574000",
			"responseCode": "000",
			"responseMessage": "SUCCESS",
			"accountName": "TEST USER",
			"ccy": "GHS",
			"branchCode": "H01",
			"customerID": "410592151",
			"availableBalance": 15.92,
			"currentBalance": 15.92,
			"odlimit": 0,
			"accountType": "S",
			"accountClass": "KEXSAV",
			"accountStatus": "ACTIVE"
		},
		"response_timestamp": "2022-04-19T19:46:57.557"
	}`

	client := newMockClient(t, mockResponse, http.StatusOK)

	opt := &AccountBalanceOptions{
		RequestID:     "14232436312",
		AffiliateCode: "EGH",
		AccountNo:     "6500184371",
		ClientID:      "ECO00184371123",
		CompanyName:   "ECOBANK TEST CO",
	}

	resp, _, err := client.Account.GetBalance(t.Context(), opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1441000574000", resp.AccountNo)
	assert.Equal(t, "TEST USER", resp.AccountName)
	assert.Equal(t, "GHS", resp.Currency)
	assert.Equal(t, decimal.NewFromFloat(15.92), resp.AvailableBalance)
	assert.Equal(t, decimal.NewFromFloat(15.92), resp.CurrentBalance)
	assert.Equal(t, "ACTIVE", resp.AccountStatus)
}

func TestAccountService_Enquiry(t *testing.T) {
	mockResponse := `{
		"response_code": 200,
		"response_message": "success",
		"response_content": {
			"accountNo": "1441000574000",
			"accountName": "TEST USER",
			"ccy": "GHS",
			"accountStatus": "ACTIVE",
			"responseCode": "000",
			"responseMessage": "SUCCESS",
			"affiliateCode": "EGH",
			"requestId": "ECO00184371123",
			"sourceCode": "CORPORATEAPI"
		},
		"response_timestamp": "2022-04-19T19:52:51.596"
	}`

	client := newMockClient(t, mockResponse, http.StatusOK)

	opt := &AccountEnquiryOptions{
		RequestID:     "14232436312",
		AffiliateCode: "EGH",
		AccountNo:     "1441000574000",
		ClientID:      "ECO00184371123",
		CompanyName:   "ECOBANK TEST CO",
	}

	resp, _, err := client.Account.Enquiry(t.Context(), opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "1441000574000", resp.AccountNo)
	assert.Equal(t, "TEST USER", resp.AccountName)
	assert.Equal(t, "GHS", resp.Currency)
	assert.Equal(t, "ACTIVE", resp.AccountStatus)
}

func TestAccountService_EnquiryThirdParty(t *testing.T) {
	mockResponse := `{
		"response_code": 200,
		"response_message": "success",
		"response_content": {
			"accountName": "PURCHASE ACCOUNT",
			"accountType": "S",
			"accountStatus": "ACTIVE",
			"hostHeaderInfo": {
				"sourceCode": "CORPORATEAPI",
				"requestId": "726262198272",
				"affiliateCode": "EGH",
				"responseCode": "000",
				"responseMessage": "success"
			}
		},
		"response_timestamp": "2021-11-03T18:20:05.058"
	}`

	client := newMockClient(t, mockResponse, http.StatusOK)

	opt := &AccountEnquiryThirdPartyOptions{
		RequestID:           "726262198272",
		AffiliateCode:       "EGH",
		AccountNo:           "1020820171412",
		DestinationBankCode: "300315",
		ClientID:            "EC06500184371123",
		CompanyName:         "Ecobanker",
	}

	resp, _, err := client.Account.EnquiryThirdParty(t.Context(), opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "PURCHASE ACCOUNT", resp.AccountName)
	assert.Equal(t, "S", resp.AccountType)
	assert.Equal(t, "ACTIVE", resp.AccountStatus)
	assert.Equal(t, "000", resp.HostHeaderInfo.ResponseCode)
}

func TestAccountService_GenerateStatement(t *testing.T) {
	mockResponse := `{
		"response_code": 200,
		"response_message": "success",
		"response_content": [
			{
				"acccy": "GHS",
				"drcrind": "CR",
				"trnrefno": "H75ZEXA1923800E0",
				"paidin": "10",
				"paidout": "",
				"valuedate": "2019-09-02 20:00:00.0",
				"lcyamount1": "10",
				"narrative": "MOBILE TRANSFER BD1441000820520-SA Xpress Account DT0209"
			},
			{
				"acccy": "GHS",
				"drcrind": "CR",
				"trnrefno": "H75ZEXA1923800E1",
				"paidin": "15",
				"paidout": "",
				"valuedate": "2019-09-02 21:30:00.0",
				"lcyamount1": "15",
				"narrative": "MOBILE TRANSFER BD1441000820520-SA Xpress Account DT0201"
			},
			{
				"acccy": "GHS",
				"drcrind": "CR",
				"trnrefno": "H75ZEXA1923800E2",
				"paidin": "17",
				"paidout": "",
				"valuedate": "2019-09-02 10:03:00.0",
				"lcyamount1": "17",
				"narrative": "MOBILE TRANSFER BD1441000820520-SA Xpress Account DT0205"
			}
		],
		"response_timestamp": "2022-04-19T19:44:21.866"
	}`

	client := newMockClient(t, mockResponse, http.StatusOK)

	opt := &GenerateStatementOptions{
		CorporateID:   "OMNI",
		RequestID:     "123456",
		ClientID:      "ZEEPAY",
		AffiliateCode: "EGH",
		AccountNumber: "1441000574000",
		StartDate:     Date{Time: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)},
		EndDate:       Date{Time: time.Date(2020, 3, 16, 0, 0, 0, 0, time.UTC)},
	}

	resp, _, err := client.Account.GenerateStatement(t.Context(), opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp, 3)
	assert.Equal(t, "H75ZEXA1923800E0", resp[0].RefNumber)
	assert.Equal(t, "10", resp[0].PaidIn)
	assert.Equal(t, "MOBILE TRANSFER BD1441000820520-SA Xpress Account DT0209", resp[0].Narrative)
}
