package contracts

type AuthActionResponse struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	RedirectURL string           `json:"redirectUrl,omitempty"`
	Session     *AuthSessionUser `json:"session,omitempty"`
}

type AuthProfileFields struct {
	PhoneNumber           string `json:"phoneNumber"`
	StoreName             string `json:"storeName"`
	SavedCardHolderName   string `json:"savedCardHolderName"`
	SavedCardNumber       string `json:"savedCardNumber"`
	SavedCardExpiry       string `json:"savedCardExpiry"`
	BankName              string `json:"bankName"`
	BankAccountNumber     string `json:"bankAccountNumber"`
	BankBookImageName     string `json:"bankBookImageName"`
	BankBookImageURL      string `json:"bankBookImageUrl"`
	IdentityCardImageName string `json:"identityCardImageName"`
	IdentityCardImageURL  string `json:"identityCardImageUrl"`
	NotifyOrders          bool   `json:"notifyOrders"`
	NotifyMarketplace     bool   `json:"notifyMarketplace"`
}

type AuthSessionUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Provider string `json:"provider"`
	IsMock   bool   `json:"isMock"`
	AuthProfileFields
}

type AuthStartRequest struct {
	Intent string `json:"intent" binding:"required,oneof=login register"`
}

type GoogleSessionExchangeRequest struct {
	AccessToken string `json:"accessToken" binding:"required"`
	Intent      string `json:"intent" binding:"omitempty,oneof=login register"`
}

type CartItem struct {
	ID         string `json:"id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Category   string `json:"category" binding:"required"`
	AuthorName string `json:"authorName" binding:"required"`
	Price      int64  `json:"price" binding:"required"`
	License    string `json:"license" binding:"required"`
}

type ProductResponse struct {
	ID                string   `json:"id"`
	CategoryID        string   `json:"categoryId"`
	Title             string   `json:"title"`
	Summary           string   `json:"summary"`
	FullDescription   string   `json:"fullDescription"`
	Category          string   `json:"category"`
	Price             int64    `json:"price"`
	Rating            float64  `json:"rating"`
	Sales             int      `json:"sales"`
	Tags              []string `json:"tags"`
	Stack             []string `json:"stack"`
	FeatureHighlights []string `json:"featureHighlights"`
	IncludedItems     []string `json:"includedItems"`
	IdealFor          []string `json:"idealFor"`
	SupportInfo       string   `json:"supportInfo"`
	VersionLabel      string   `json:"versionLabel"`
	FileFormatLabel   string   `json:"fileFormatLabel"`
	AuthorName        string   `json:"authorName"`
	AuthorSlug        string   `json:"authorSlug"`
	UpdatedAt         string   `json:"updatedAt"`
	UpdatedDaysAgo    int      `json:"updatedDaysAgo"`
	Delivery          string   `json:"delivery"`
	License           string   `json:"license"`
	LicenseID         string   `json:"licenseId"`
	Verified          bool     `json:"verified"`
}

type MarketplaceSellerResponse struct {
	Slug                string   `json:"slug"`
	Name                string   `json:"name"`
	Summary             string   `json:"summary"`
	ProductCount        int      `json:"productCount"`
	TotalSales          int      `json:"totalSales"`
	StartingPrice       int64    `json:"startingPrice"`
	VerifiedCount       int      `json:"verifiedCount"`
	Categories          []string `json:"categories"`
	Stacks              []string `json:"stacks"`
	LatestUpdateDaysAgo int      `json:"latestUpdateDaysAgo"`
}

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

type CheckoutSubmitInput struct {
	CustomerName           string                 `json:"customerName" binding:"required"`
	CustomerEmail          string                 `json:"customerEmail" binding:"required,email"`
	CustomerPhone          string                 `json:"customerPhone" binding:"required"`
	CompanyName            string                 `json:"companyName"`
	TaxID                  string                 `json:"taxId"`
	Note                   string                 `json:"note"`
	PaymentMethod          string                 `json:"paymentMethod" binding:"required,oneof=promptpay card bank-transfer"`
	ReceivePurchaseUpdates bool                   `json:"receivePurchaseUpdates"`
	RequireInvoice         bool                   `json:"requireInvoice"`
	Subtotal               int64                  `json:"subtotal" binding:"required"`
	Total                  int64                  `json:"total" binding:"required"`
	Items                  []CartItem             `json:"items" binding:"required,min=1"`
	PaymentDetails         CheckoutPaymentDetails `json:"paymentDetails"`
}

type CheckoutSubmitResponse struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	OrderID       string `json:"orderId"`
	PaymentMethod string `json:"paymentMethod"`
	Total         int64  `json:"total"`
	Status        string `json:"status"`
}

type SellerListingRequest struct {
	AssetType             string   `json:"assetType" binding:"required,oneof=source-code template component-kit"`
	Title                 string   `json:"title" binding:"required"`
	CategoryID            string   `json:"categoryId" binding:"required"`
	LicenseID             string   `json:"licenseId" binding:"required"`
	Price                 int64    `json:"price" binding:"required"`
	Summary               string   `json:"summary" binding:"required"`
	Description           string   `json:"description" binding:"required"`
	Highlights            []string `json:"highlights"`
	IdealFor              []string `json:"idealFor"`
	SupportInfo           string   `json:"supportInfo"`
	Stack                 []string `json:"stack"`
	Version               string   `json:"version"`
	DemoURL               string   `json:"demoUrl"`
	SupportURL            string   `json:"supportUrl"`
	IncludedFiles         []string `json:"includedFiles"`
	PackageFileName       string   `json:"packageFileName"`
	CoverFileName         string   `json:"coverFileName"`
	DocsFileName          string   `json:"docsFileName"`
	InstantDelivery       bool     `json:"instantDelivery"`
	SourceIncluded        bool     `json:"sourceIncluded"`
	DocumentationIncluded bool     `json:"documentationIncluded"`
	Mode                  string   `json:"mode" binding:"required,oneof=draft publish"`
}

type SellerListingResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ListingID   string `json:"listingId"`
	Status      string `json:"status"`
}

type SellerOrderResponse struct {
	ID                 string `json:"id"`
	OrderID            string `json:"orderId"`
	ProductID          string `json:"productId"`
	ProductTitle       string `json:"productTitle"`
	ProductCategory    string `json:"productCategory"`
	BuyerName          string `json:"buyerName"`
	BuyerEmail         string `json:"buyerEmail"`
	PurchasedAt        string `json:"purchasedAt"`
	Amount             int64  `json:"amount"`
	PaymentMethodLabel string `json:"paymentMethodLabel"`
	LicenseLabel       string `json:"licenseLabel"`
	DeliveryLabel      string `json:"deliveryLabel"`
	StatusLabel        string `json:"statusLabel"`
}

type DownloadLibraryItemResponse struct {
	LibraryItemID    string  `json:"libraryItemId"`
	OrderID          string  `json:"orderId"`
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Category         string  `json:"category"`
	AuthorName       string  `json:"authorName"`
	Price            int64   `json:"price"`
	License          string  `json:"license"`
	PurchasedAt      string  `json:"purchasedAt"`
	PaymentMethod    string  `json:"paymentMethod"`
	Status           string  `json:"status"`
	VersionLabel     string  `json:"versionLabel"`
	FileName         string  `json:"fileName"`
	FileSizeLabel    string  `json:"fileSizeLabel"`
	DownloadsCount   int     `json:"downloadsCount"`
	LastDownloadedAt *string `json:"lastDownloadedAt"`
}

type CookiePreferences struct {
	Necessary   bool `json:"necessary"`
	Preferences bool `json:"preferences"`
	Analytics   bool `json:"analytics"`
	Marketing   bool `json:"marketing"`
}

type CookieConsentResponse struct {
	Status      string            `json:"status"`
	Preferences CookiePreferences `json:"preferences"`
	UpdatedAt   string            `json:"updatedAt"`
}

type CookieConsentUpsertRequest struct {
	Status      string            `json:"status" binding:"required,oneof=necessary all customized"`
	Preferences CookiePreferences `json:"preferences"`
}

type ProfileUpdateRequest struct {
	Name                  string `json:"name"`
	PhoneNumber           string `json:"phoneNumber"`
	StoreName             string `json:"storeName"`
	SavedCardHolderName   string `json:"savedCardHolderName"`
	SavedCardNumber       string `json:"savedCardNumber"`
	SavedCardExpiry       string `json:"savedCardExpiry"`
	BankName              string `json:"bankName"`
	BankAccountNumber     string `json:"bankAccountNumber"`
	BankBookImageName     string `json:"bankBookImageName"`
	BankBookImageURL      string `json:"bankBookImageUrl"`
	IdentityCardImageName string `json:"identityCardImageName"`
	IdentityCardImageURL  string `json:"identityCardImageUrl"`
	NotifyOrders          bool   `json:"notifyOrders"`
	NotifyMarketplace     bool   `json:"notifyMarketplace"`
}

type MessageResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
