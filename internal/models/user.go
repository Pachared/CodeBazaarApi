package models

type User struct {
	ID                    string    `gorm:"primaryKey;size:64"`
	Slug                  string    `gorm:"size:120;index"`
	Name                  string    `gorm:"size:255;not null"`
	Email                 string    `gorm:"size:255;uniqueIndex;not null"`
	Role                  string    `gorm:"size:32;index;not null"`
	Provider              string    `gorm:"size:32;not null;default:google"`
	PhoneNumber           string    `gorm:"size:64"`
	StoreName             string    `gorm:"size:255"`
	SavedCardHolderName   string    `gorm:"size:255"`
	SavedCardNumber       string    `gorm:"size:64"`
	SavedCardExpiry       string    `gorm:"size:16"`
	BankName              string    `gorm:"size:255"`
	BankAccountNumber     string    `gorm:"size:64"`
	BankBookImageName     string    `gorm:"size:255"`
	BankBookImageURL      string    `gorm:"type:text"`
	IdentityCardImageName string    `gorm:"size:255"`
	IdentityCardImageURL  string    `gorm:"type:text"`
	NotifyOrders          bool      `gorm:"default:true"`
	NotifyMarketplace     bool      `gorm:"default:true"`
	Products              []Product `gorm:"foreignKey:SellerID"`
	TimestampedModel
}
