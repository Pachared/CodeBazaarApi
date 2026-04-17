package services

import (
	"errors"
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/models"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
	"gorm.io/gorm"
)

type CookieService struct {
	cookieConsentRepository *repositories.CookieConsentRepository
}

func NewCookieService(cookieConsentRepository *repositories.CookieConsentRepository) *CookieService {
	return &CookieService{cookieConsentRepository: cookieConsentRepository}
}

func (s *CookieService) GetConsent(currentUser *models.User, sessionKey string) (*contracts.CookieConsentResponse, error) {
	userID := ""
	if currentUser != nil {
		userID = currentUser.ID
	}

	consent, err := s.cookieConsentRepository.GetByUserIDOrSession(userID, sessionKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return toCookieConsentResponse(consent), nil
}

func (s *CookieService) SaveConsent(
	currentUser *models.User,
	sessionKey string,
	input contracts.CookieConsentUpsertRequest,
) (*contracts.CookieConsentResponse, error) {
	userID := ""
	if currentUser != nil {
		userID = currentUser.ID
	}
	if userID == "" && sessionKey == "" {
		sessionKey = "anonymous-demo"
	}

	existing, _ := s.cookieConsentRepository.GetByUserIDOrSession(userID, sessionKey)
	if existing == nil {
		existing = &models.CookieConsent{
			ID:         createStableID("cookie"),
			UserID:     userID,
			SessionKey: sessionKey,
		}
	}

	existing.UserID = userID
	existing.SessionKey = sessionKey
	existing.Status = input.Status
	existing.Preferences = models.CookiePreferences{
		Necessary:   true,
		Preferences: input.Preferences.Preferences,
		Analytics:   input.Preferences.Analytics,
		Marketing:   input.Preferences.Marketing,
	}
	existing.ConsentUpdatedAt = time.Now()

	if err := s.cookieConsentRepository.Save(existing); err != nil {
		return nil, err
	}

	return toCookieConsentResponse(existing), nil
}
