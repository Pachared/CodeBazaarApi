package models

import "time"

type DownloadItem struct {
	ID               string    `gorm:"primaryKey;size:64"`
	UserID           string    `gorm:"size:64;index;not null"`
	OrderID          string    `gorm:"size:64;index;not null"`
	OrderItemID      string    `gorm:"size:64;uniqueIndex;not null"`
	ProductID        string    `gorm:"size:120;index;not null"`
	Title            string    `gorm:"size:255;not null"`
	Category         string    `gorm:"size:120;not null"`
	AuthorName       string    `gorm:"size:255;not null"`
	Price            int64     `gorm:"not null"`
	License          string    `gorm:"size:120"`
	PaymentMethod    string    `gorm:"size:64;not null"`
	PurchasedAt      time.Time `gorm:"index;not null"`
	Status           string    `gorm:"size:32;not null;default:ready"`
	VersionLabel     string    `gorm:"size:64"`
	FileName         string    `gorm:"size:255;not null"`
	FileSizeLabel    string    `gorm:"size:64;not null"`
	DownloadsCount   int       `gorm:"default:0"`
	LastDownloadedAt *time.Time
	TimestampedModel
}
