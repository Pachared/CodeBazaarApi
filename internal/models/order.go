package models

import "time"

type Order struct {
	ID                     string                 `gorm:"primaryKey;size:64"`
	OrderNumber            string                 `gorm:"size:64;uniqueIndex;not null"`
	BuyerID                string                 `gorm:"size:64;index"`
	Buyer                  User                   `gorm:"foreignKey:BuyerID"`
	CustomerName           string                 `gorm:"size:255;not null"`
	CustomerEmail          string                 `gorm:"size:255;not null"`
	CustomerPhone          string                 `gorm:"size:64;not null"`
	CompanyName            string                 `gorm:"size:255"`
	TaxID                  string                 `gorm:"size:64"`
	Note                   string                 `gorm:"type:text"`
	PaymentMethod          string                 `gorm:"size:64;not null"`
	ReceivePurchaseUpdates bool                   `gorm:"default:true"`
	RequireInvoice         bool                   `gorm:"default:false"`
	Subtotal               int64                  `gorm:"not null"`
	Total                  int64                  `gorm:"not null"`
	Status                 string                 `gorm:"size:32;not null;default:paid"`
	PaymentDetails         CheckoutPaymentDetails `gorm:"type:jsonb;serializer:json"`
	PurchasedAt            time.Time              `gorm:"index;not null"`
	Items                  []OrderItem            `gorm:"foreignKey:OrderID"`
	TimestampedModel
}

type OrderItem struct {
	ID            string `gorm:"primaryKey;size:64"`
	OrderID       string `gorm:"size:64;index;not null"`
	ProductID     string `gorm:"size:120;index;not null"`
	SellerID      string `gorm:"size:64;index;not null"`
	Title         string `gorm:"size:255;not null"`
	Category      string `gorm:"size:120;not null"`
	AuthorName    string `gorm:"size:255;not null"`
	Price         int64  `gorm:"not null"`
	License       string `gorm:"size:120"`
	DeliveryLabel string `gorm:"size:255"`
	TimestampedModel
}
