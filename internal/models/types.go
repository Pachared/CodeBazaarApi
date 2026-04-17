package models

type PromptPayPaymentDetails struct {
	AccountName   string `json:"accountName"`
	PromptPayID   string `json:"promptPayId"`
	ReferenceCode string `json:"referenceCode"`
}

type CardPaymentDetails struct {
	CardHolderName   string `json:"cardHolderName"`
	CardNumberMasked string `json:"cardNumberMasked"`
	Expiry           string `json:"expiry"`
}

type BankTransferPaymentDetails struct {
	BankName          string `json:"bankName"`
	AccountName       string `json:"accountName"`
	AccountNumber     string `json:"accountNumber"`
	TransferReference string `json:"transferReference"`
	SlipImageName     string `json:"slipImageName"`
}

type CheckoutPaymentDetails struct {
	PromptPay    *PromptPayPaymentDetails    `json:"promptpay,omitempty"`
	Card         *CardPaymentDetails         `json:"card,omitempty"`
	BankTransfer *BankTransferPaymentDetails `json:"bankTransfer,omitempty"`
}

type CookiePreferences struct {
	Necessary   bool `json:"necessary"`
	Preferences bool `json:"preferences"`
	Analytics   bool `json:"analytics"`
	Marketing   bool `json:"marketing"`
}
