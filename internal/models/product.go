package models

type Product struct {
	ID                    string   `gorm:"primaryKey;size:120"`
	SellerID              string   `gorm:"size:64;index;not null"`
	Seller                User     `gorm:"foreignKey:SellerID"`
	AssetType             string   `gorm:"size:64;not null;default:source-code"`
	Status                string   `gorm:"size:32;index;not null;default:publish"`
	CategoryID            string   `gorm:"size:64;index;not null"`
	Category              string   `gorm:"size:120;not null"`
	Title                 string   `gorm:"size:255;not null"`
	Summary               string   `gorm:"type:text;not null"`
	FullDescription       string   `gorm:"type:text;not null"`
	Price                 int64    `gorm:"not null"`
	Rating                float64  `gorm:"type:numeric(3,2);default:0"`
	Sales                 int      `gorm:"default:0"`
	Tags                  []string `gorm:"type:jsonb;serializer:json"`
	Stack                 []string `gorm:"type:jsonb;serializer:json"`
	FeatureHighlights     []string `gorm:"type:jsonb;serializer:json"`
	IncludedItems         []string `gorm:"type:jsonb;serializer:json"`
	IdealFor              []string `gorm:"type:jsonb;serializer:json"`
	SupportInfo           string   `gorm:"type:text"`
	VersionLabel          string   `gorm:"size:64"`
	FileFormatLabel       string   `gorm:"size:255"`
	UpdatedDaysAgo        int      `gorm:"default:0"`
	Delivery              string   `gorm:"size:255"`
	License               string   `gorm:"size:120"`
	LicenseID             string   `gorm:"size:64;index"`
	Verified              bool     `gorm:"default:false"`
	DemoURL               string   `gorm:"type:text"`
	SupportURL            string   `gorm:"type:text"`
	PackageFileName       string   `gorm:"size:255"`
	CoverFileName         string   `gorm:"size:255"`
	DocsFileName          string   `gorm:"size:255"`
	InstantDelivery       bool     `gorm:"default:true"`
	SourceIncluded        bool     `gorm:"default:true"`
	DocumentationIncluded bool     `gorm:"default:true"`
	TimestampedModel
}
