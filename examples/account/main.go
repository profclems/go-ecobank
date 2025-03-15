package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/profclems/go-ecobank"
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

	// get account balance
	acctBal, resp, err := client.Account.GetBalance(ctx, &ecobank.AccountBalanceOptions{
		RequestID:     "14232436312",
		AffiliateCode: "EGH",
		AccountNo:     "6500184371",
		ClientID:      "ECO00184371123",
		CompanyName:   "ECOBANK TEST CO",
	})
	checkErr(errors.Wrap(err, "failed to get balance"))

	fmt.Println("Code:", resp.Code)
	fmt.Println("Message:", resp.Message)
	checkErr(json.NewEncoder(os.Stdout).Encode(acctBal))
	fmt.Println()

	// get account details
	enquiry, resp, err := client.Account.Enquiry(ctx, &ecobank.AccountEnquiryOptions{
		RequestID:     "14232436312",
		AffiliateCode: "EGH",
		AccountNo:     "1441000574000",
		ClientID:      "ECO00184371123",
		CompanyName:   "ECOBANK TEST CO",
	})
	checkErr(errors.Wrap(err, "failed to get account details"))

	fmt.Println("Code:", resp.Code)
	fmt.Println("Message:", resp.Message)
	checkErr(json.NewEncoder(os.Stdout).Encode(enquiry))
	fmt.Println()

	enquiryTP, resp, err := client.Account.EnquiryThirdParty(ctx, &ecobank.AccountEnquiryThirdPartyOptions{
		RequestID:           "726262198272",
		AffiliateCode:       "EGH",
		AccountNo:           "1020820171412",
		DestinationBankCode: "300315",
		ClientID:            "EC06500184371123",
		CompanyName:         "Ecobanker",
	})
	checkErr(errors.Wrap(err, "failed to get third party account details"))
	fmt.Println("Code:", resp.Code)
	fmt.Println("Message:", resp.Message)
	checkErr(json.NewEncoder(os.Stdout).Encode(enquiryTP))
	fmt.Println()

	statementOptions := &ecobank.GenerateStatementOptions{
		RequestID:     "123456",
		ClientID:      "ZEEPAY",
		AffiliateCode: "EGH",
		CorporateID:   "OMNI",
		AccountNumber: "1441000574000",
		StartDate:     ecobank.NewDate(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:       ecobank.NewDate(time.Date(2020, 3, 16, 0, 0, 0, 0, time.UTC)),
	}
	// generate account statement
	statement, resp, err := client.Account.GenerateStatement(ctx, statementOptions)
	checkErr(errors.Wrap(err, "failed to generate statement"))

	fmt.Println("Code:", resp.Code)
	fmt.Println("Message:", resp.Message)
	checkErr(json.NewEncoder(os.Stdout).Encode(statement))
	fmt.Println()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
