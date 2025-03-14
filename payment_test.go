package ecobank

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentService_GetBillerList(t *testing.T) {
	mockResponse := `{
		"response_code": 200,
		"response_message": "success",
		"response_content": {
			"hostHeaderInfo": {
				"sourceCode": "ECOBANKMOBILEAPP",
				"requestId": "ECO2112134345",
				"affiliateCode": "EGH",
				"responseCode": "000",
				"responseMessage": "Success"
			},
			"billerInfo": [
				{
					"billerCode": "MGC",
					"billerID": 77427,
					"billerName": "METHODIST COLLECTION",
					"billerDescription": "METHODIST COLLECTION",
					"billerCategory": null,
					"billerLogo": "/usr/app/Alert/ecobank_banner.jpg",
					"billAmountType": "",
					"billAmount": 0,
					"ccy": "",
					"collectionAccountNo": "",
					"aggregatorName": "NEWESB",
					"amountDenominations": "",
					"productCodeList": ""
				},
				{
					"billerCode": "GHWATER",
					"billerID": 76758,
					"billerName": "GHANA WATER",
					"billerDescription": "GHANA WATER",
					"billerCategory": "ECOBANK",
					"billerLogo": "/usr/app/Alert/ecobank_banner.jpg",
					"billAmountType": "",
					"billAmount": 1,
					"ccy": "GHS",
					"collectionAccountNo": "",
					"aggregatorName": "GHANA WATER",
					"amountDenominations": "",
					"productCodeList": ""
				}
			]
		},
		"response_timestamp": "2022-09-23T17:04:43.506"
	}`

	client := newMockClient(t, mockResponse, http.StatusOK)

	opt := &GetBillerListOptions{
		RequestID:     "ECO2112134345",
		AffiliateCode: "EGH",
	}

	resp, _, err := client.Payment.GetBillerList(t.Context(), opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.BillerInfo, 2)

	// Validate first biller
	assert.Equal(t, "MGC", resp.BillerInfo[0].BillerCode)
	assert.Equal(t, 77427, resp.BillerInfo[0].BillerID)
	assert.Equal(t, "METHODIST COLLECTION", resp.BillerInfo[0].BillerName)
	assert.Equal(t, "/usr/app/Alert/ecobank_banner.jpg", resp.BillerInfo[0].BillerLogo)
	assert.Equal(t, "NEWESB", resp.BillerInfo[0].AggregatorName)
	assert.Equal(t, 0, resp.BillerInfo[0].BillAmount)

	// Validate second biller
	assert.Equal(t, "GHWATER", resp.BillerInfo[1].BillerCode)
	assert.Equal(t, 76758, resp.BillerInfo[1].BillerID)
	assert.Equal(t, "GHANA WATER", resp.BillerInfo[1].BillerName)
	assert.Equal(t, "ECOBANK", resp.BillerInfo[1].BillerCategory)
	assert.Equal(t, "/usr/app/Alert/ecobank_banner.jpg", resp.BillerInfo[1].BillerLogo)
	assert.Equal(t, "GHANA WATER", resp.BillerInfo[1].AggregatorName)
	assert.Equal(t, 1, resp.BillerInfo[1].BillAmount)
	assert.Equal(t, "GHS", resp.BillerInfo[1].Currency)

	// Validate host header info
	assert.Equal(t, "ECOBANKMOBILEAPP", resp.HostHeaderInfo.SourceCode)
	assert.Equal(t, "ECO2112134345", resp.HostHeaderInfo.RequestID)
	assert.Equal(t, "EGH", resp.HostHeaderInfo.AffiliateCode)
	assert.Equal(t, "000", resp.HostHeaderInfo.ResponseCode)
	assert.Equal(t, "Success", resp.HostHeaderInfo.ResponseMessage)
}

func TestPaymentService_ValidateBiller(t *testing.T) {
	mockResponse := `{
 "response_code": 200,
 "response_message": "success",
 "response_content": {
  "hostHeaderInfo": {
   "sourceCode": "ECOBANKMOBILEAPP",
   "requestId": "0254875943",
   "affiliateCode": "EGH",
   "responseCode": "000",
   "responseMessage": "Success"
  },
  "billerCode": "MTNPTU",
  "billRefNo": "46356262",
  "customerName": "Benson",
  "amount": 0,
  "paymentDescription": "",
  "productCode": "",
  "responseValues": "",
  "formDataValue": [
   {
    "fieldName": "CHARGE",
    "fieldDescription": "",
    "fieldMasked": "",
    "fieldValue": "100.0",
    "fieldRequired": "",
    "dataType": "DOUBLE"
   },
   {
    "fieldName": "VAT",
    "fieldDescription": "",
    "fieldMasked": "",
    "fieldValue": "0.0",
    "fieldRequired": "",
    "dataType": "DOUBLE"
   },
   {
    "fieldName": "VAT",
    "fieldDescription": "",
    "fieldMasked": "",
    "fieldValue": "0.0",
    "fieldRequired": "",
    "dataType": "DOUBLE"
   },
   {
    "fieldName": "TOTAL CHARGE",
    "fieldDescription": "",
    "fieldMasked": "",
    "fieldValue": "100.0",
    "fieldRequired": "",
    "dataType": "DOUBLE"
   },
   {
    "fieldName": "TOTAL CHARGE",
    "fieldDescription": "",
    "fieldMasked": "",
    "fieldValue": "100.0",
    "fieldRequired": "",
    "dataType": "DOUBLE"
   }
  ]
 },
 "response_timestamp": "2022-09-23T17:17:53.181"
}`

	client := newMockClient(t, mockResponse, http.StatusOK)

	opt := &ValidateBillerOptions{
		RequestID:     "EC12O2134521",
		AffiliateCode: "EGH",
		BillerCode:    "MTNPTU",
		ProductCode:   "02",
		MobileNumber:  "0254875943",
		CustomerName:  "Edu",
		FormDataValue: []struct {
			FieldName  string `json:"fieldName"`
			FieldValue string `json:"fieldValue"`
		}{{
			FieldName:  "METER NUMBER",
			FieldValue: "54140081982",
		}},
	}

	resp, _, err := client.Payment.ValidateBiller(t.Context(), opt)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, "MTNPTU", resp.BillerCode)
	assert.Equal(t, "46356262", resp.BillRefNo)
	assert.Equal(t, "Benson", resp.CustomerName)
	assert.Equal(t, 0, resp.Amount)
	assert.Equal(t, "", resp.PaymentDescription)
	assert.Equal(t, "", resp.ProductCode)
	assert.Equal(t, "", resp.ResponseValues)

	assert.Len(t, resp.FormDataValue, 5)

	assert.Equal(t, "CHARGE", resp.FormDataValue[0].FieldName)
	assert.Equal(t, "100.0", resp.FormDataValue[0].FieldValue)
	assert.Equal(t, "DOUBLE", resp.FormDataValue[0].DataType)

	assert.Equal(t, "VAT", resp.FormDataValue[1].FieldName)
	assert.Equal(t, "0.0", resp.FormDataValue[1].FieldValue)
	assert.Equal(t, "DOUBLE", resp.FormDataValue[1].DataType)

	assert.Equal(t, "TOTAL CHARGE", resp.FormDataValue[3].FieldName)
	assert.Equal(t, "100.0", resp.FormDataValue[3].FieldValue)
	assert.Equal(t, "DOUBLE", resp.FormDataValue[3].DataType)

	assert.Equal(t, "ECOBANKMOBILEAPP", resp.HostHeaderInfo.SourceCode)
	assert.Equal(t, "0254875943", resp.HostHeaderInfo.RequestID)
	assert.Equal(t, "EGH", resp.HostHeaderInfo.AffiliateCode)
	assert.Equal(t, "000", resp.HostHeaderInfo.ResponseCode)
	assert.Equal(t, "Success", resp.HostHeaderInfo.ResponseMessage)
}
