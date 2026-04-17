package repositories

import (
	"errors"

	"github.com/Pachared/CodeBazaarApi/internal/models"
	"gorm.io/gorm"
)

type CookieConsentRepository struct {
	db *gorm.DB
}

func NewCookieConsentRepository(db *gorm.DB) *CookieConsentRepository {
	return &CookieConsentRepository{db: db}
}

func (r *CookieConsentRepository) GetByUserIDOrSession(userID string, sessionKey string) (*models.CookieConsent, error) {
	var consent models.CookieConsent

	query := r.db.Model(&models.CookieConsent{})
	switch {
	case userID != "":
		query = query.Where("user_id = ?", userID)
	case sessionKey != "":
		query = query.Where("session_key = ?", sessionKey)
	default:
		return nil, gorm.ErrRecordNotFound
	}

	if err := query.Order("consent_updated_at DESC").First(&consent).Error; err != nil {
		return nil, err
	}

	return &consent, nil
}

func (r *CookieConsentRepository) Save(consent *models.CookieConsent) error {
	if consent == nil {
		return errors.New("cookie consent is nil")
	}

	return r.db.Save(consent).Error
}
