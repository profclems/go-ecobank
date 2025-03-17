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

	fmt.Println("Creating account...")
	// open account
	createOpts := &ecobank.CreateAccountOptions{
		ClientID:           "ECO76383823",
		RequestID:          "ECO76383823",
		AffiliateCode:      "ENG",
		FirstName:          "Rotimi",
		Middlename:         "",
		Lastname:           "Akinola",
		MobileNo:           "2348089991325",
		Gender:             "M",
		IdentityNo:         "198837383982",
		IdentityType:       "MOBILE_WALLET_NO",
		IDIssueDate:        "01072021",
		IDExpiryDate:       "01072021",
		Ccy:                "NGN",
		Country:            "NGN",
		BranchCode:         "ENG",
		DateOfBirth:        "01072021",
		CountryOfResidence: "NIGERIA",
		Email:              "treknfreedom@yahoo.com",
		Street:             "Labone",
		City:               "Accra",
		State:              "Accra",
		Image:              "oeyetweuiww8262822999999999",
		Signature:          "orjerjeklellwewpw726527289292",
	}
	createOpts.SetHash("a43aa74662060b7b9c942dd7ace565a0919118db758bcd71a0f5c7cd7e349f6309b02866b6156ef9171a1b23119c71e77db2edd38cc89963d7f34b541d6dc461")
	account, resp, err := client.Account.CreateAccount(ctx, createOpts)
	checkErr(errors.Wrap(err, "failed to create account"))

	fmt.Println("Code:", resp.Code)
	fmt.Println("Message:", resp.Message)
	fmt.Println("Account Number:", account.AccountNo)
	fmt.Println("Client ID:", account.ClientID)
	fmt.Println("Track Ref:", account.TrackRef)
	fmt.Println("Short Name:", account.Shortname)
	fmt.Println("Mobile Number:", account.MobileNo)
	fmt.Println()

	// get account balance
	fmt.Println("Getting account balance...")
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
	fmt.Println("Account Number:", acctBal.AccountNo)
	fmt.Println("Account Name:", acctBal.AccountName)
	fmt.Println("Available Balance:", acctBal.AvailableBalance)
	fmt.Println("Current Balance:", acctBal.CurrentBalance)
	fmt.Println("Currency:", acctBal.Currency)
	fmt.Println()

	// get account details
	fmt.Println("Getting account details...")
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
	fmt.Println("Account Name:", enquiry.AccountName)
	fmt.Println("Account Number:", enquiry.AccountNo)
	fmt.Println("Account Status:", enquiry.AccountStatus)
	fmt.Println("Currency:", enquiry.Currency)
	fmt.Println()

	fmt.Println("Getting third party account details")
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
	fmt.Println("Account Name:", enquiryTP.AccountName)
	fmt.Println("Account Number:", enquiryTP.AccountType)
	fmt.Println("Account Status:", enquiryTP.AccountStatus)
	fmt.Println()

	fmt.Println("Generating account statement...")
	statementOptions := &ecobank.GenerateStatementOptions{
		RequestID:     "123456",
		ClientID:      "ZEEPAY",
		AffiliateCode: "EGH",
		CorporateID:   "OMNI",
		AccountNumber: "1441000574000",
		StartDate:     ecobank.NewDate(time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:       ecobank.NewDate(time.Date(2020, 3, 16, 0, 0, 0, 0, time.UTC)),
	}
	statementOptions.SetHash("aa708d5f5434bc385d9b096ff663bd19abb07658e0c7c3b0580a616dded6e05218ebdef8b1c1547446993b99a04f7e65885ca44b5dc6548acbbfd2b5d1117e5c")
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
