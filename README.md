# go-ecobank

An Ecobank API Client providing a convenient way to interact with the Ecobank Corporate API.

## Coverage
This API client package covers the following services:

* **Authentication:** Handles the login process to obtain an access token.
* **Account Services:**
    * Get account balance
    * Retrieve account details
    * Perform third-party account inquiry
    * Generate account statements
    * Create new accounts
* **Payment Services:**
    * Retrieve a list of billers
    * Get details for a specific biller
    * Validate biller information
    * Initiate various payment types (Bill Payment, Token Transfer, Domestic Transfer, Interbank Transfer, Airtime Top-up, Mobile Money Transfer)
* **Transaction Status Services:**
    * Retrieve transaction status
    * Retrieve E-Token status
* **Remittance Services:**
    * Initiate various payment types
      * Cross-Border Ecobank-to-Ecobank
      * Cross-Border Interbank
      * Cross-Border Xpress Cash Token
      * Cross-Border Bank-to-Wallet (MoMo)
    * Name Enquiry
    * Institution List

## Installation

```bash
go get github.com/profclems/go-ecobank
```

## Usage

```go
import "github.com/profclems/go-ecobank"
```
Construct a new Ecobank client, obtain an access token, and start making requests.

```go
client, err := ecobank.NewClient("username", "password", "lab-key")
if err != nil {
    log.Fatal(err)
}

// login
err = client.Login()
if err != nil {
    log.Fatal(err)
}

// get account balance
acctBal, resp, err := client.Account.GetBalance(ctx, &ecobank.AccountBalanceOptions{
    RequestID:     "14232436312",
    AffiliateCode: "EGH",
    AccountNo:     "6500184371",
    ClientID:      "ECO00184371123",
    CompanyName:   "ECOBANK TEST CO",
})
```

## TODO
This library is still a work in progress as it was built based on the sandbox environment.

Also, the biggest thing this package needs is tests. I will be adding tests in the future.
