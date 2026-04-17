package models

import "time"

type CookieConsent struct {
	ID               string            `gorm:"primaryKey;size:64"`
	UserID           string            `gorm:"size:64;index"`
	SessionKey       string            `gorm:"size:128;index"`
	Status           string            `gorm:"size:32;not null"`
	Preferences      CookiePreferences `gorm:"type:jsonb;serializer:json"`
	ConsentUpdatedAt time.Time         `gorm:"index;not null"`
	TimestampedModel
}
