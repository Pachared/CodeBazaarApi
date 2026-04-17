package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", strings.TrimSpace(id)).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", strings.TrimSpace(strings.ToLower(email))).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByIDOrEmail(id string, email string) (*models.User, error) {
	if strings.TrimSpace(id) != "" {
		user, err := r.GetByID(id)
		if err == nil {
			return user, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if strings.TrimSpace(email) == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return r.GetByEmail(email)
}

func (r *UserRepository) Save(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) FindOrCreateExternalUser(
	userID string,
	email string,
	name string,
	provider string,
	role string,
) (*models.User, error) {
	normalizedID := strings.TrimSpace(userID)
	normalizedEmail := strings.TrimSpace(strings.ToLower(email))
	normalizedProvider := resolveProvider(provider, normalizedID)
	normalizedRole := resolveRole(role, normalizedProvider)

	if normalizedID == "" && normalizedEmail == "" {
		return nil, gorm.ErrRecordNotFound
	}

	existingUser, err := r.GetByIDOrEmail(normalizedID, normalizedEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUser != nil {
		existingUser.Name = coalesceString(strings.TrimSpace(name), existingUser.Name)
		existingUser.Email = coalesceString(normalizedEmail, existingUser.Email)
		existingUser.Provider = mergeProvider(existingUser.Provider, normalizedProvider)
		existingUser.Role = mergeRole(existingUser.Role, normalizedRole)

		if existingUser.Role == "seller" {
			if strings.TrimSpace(existingUser.StoreName) == "" {
				existingUser.StoreName = existingUser.Name
			}
			if strings.TrimSpace(existingUser.Slug) == "" {
				slug, slugErr := r.uniqueSlug(existingUser.Name, existingUser.ID)
				if slugErr != nil {
					return nil, slugErr
				}
				existingUser.Slug = slug
			}
		}

		if err := r.db.Save(existingUser).Error; err != nil {
			return nil, err
		}

		return existingUser, nil
	}

	displayName := strings.TrimSpace(name)
	if displayName == "" {
		displayName = deriveNameFromEmail(normalizedEmail, normalizedRole)
	}

	user := &models.User{
		ID:                coalesceString(normalizedID, createRepositoryID("usr_external")),
		Name:              displayName,
		Email:             normalizedEmail,
		Role:              normalizedRole,
		Provider:          normalizedProvider,
		NotifyOrders:      true,
		NotifyMarketplace: true,
	}

	if user.Role == "seller" {
		slug, slugErr := r.uniqueSlug(displayName, "")
		if slugErr != nil {
			return nil, slugErr
		}
		user.Slug = slug
		user.StoreName = displayName
	}

	return r.firstOrCreateByEmail(user)
}

func (r *UserRepository) FindOrCreateBuyerByEmail(name string, email string, phone string) (*models.User, error) {
	normalizedEmail := strings.TrimSpace(strings.ToLower(email))
	if normalizedEmail == "" {
		return nil, errors.New("buyer email is required")
	}

	existingUser, err := r.GetByEmail(normalizedEmail)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUser != nil {
		existingUser.Name = coalesceString(strings.TrimSpace(name), existingUser.Name)
		existingUser.PhoneNumber = coalesceString(strings.TrimSpace(phone), existingUser.PhoneNumber)
		existingUser.Provider = mergeProvider(existingUser.Provider, "guest")
		existingUser.Role = mergeRole(existingUser.Role, "buyer")

		if err := r.db.Save(existingUser).Error; err != nil {
			return nil, err
		}

		return existingUser, nil
	}

	user := &models.User{
		ID:                createRepositoryID("usr_buyer"),
		Name:              coalesceString(strings.TrimSpace(name), "ผู้ซื้อ CodeBazaar"),
		Email:             normalizedEmail,
		Role:              "buyer",
		Provider:          "guest",
		PhoneNumber:       strings.TrimSpace(phone),
		NotifyOrders:      true,
		NotifyMarketplace: true,
	}

	return r.dbSaveAndReturn(user)
}

func (r *UserRepository) EnsureSellerAccount(currentUser *models.User) (*models.User, error) {
	if currentUser == nil || strings.TrimSpace(currentUser.ID) == "" {
		return nil, gorm.ErrRecordNotFound
	}

	user, err := r.GetByIDOrEmail(currentUser.ID, currentUser.Email)
	if err != nil {
		return nil, err
	}

	user.Role = "seller"
	if strings.TrimSpace(user.StoreName) == "" {
		user.StoreName = user.Name
	}
	if strings.TrimSpace(user.Slug) == "" {
		slug, slugErr := r.uniqueSlug(user.Name, user.ID)
		if slugErr != nil {
			return nil, slugErr
		}
		user.Slug = slug
	}

	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) firstOrCreateByEmail(user *models.User) (*models.User, error) {
	normalizedEmail := strings.TrimSpace(strings.ToLower(user.Email))
	if normalizedEmail == "" {
		return nil, errors.New("user email is required")
	}

	var existing models.User
	err := r.db.First(&existing, "email = ?", normalizedEmail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user.Email = normalizedEmail
			return r.dbSaveAndReturn(user)
		}

		return nil, err
	}

	existing.Name = coalesceString(user.Name, existing.Name)
	existing.Email = normalizedEmail
	existing.Role = mergeRole(existing.Role, user.Role)
	existing.Provider = mergeProvider(existing.Provider, user.Provider)
	existing.PhoneNumber = coalesceString(user.PhoneNumber, existing.PhoneNumber)
	existing.StoreName = coalesceString(user.StoreName, existing.StoreName)
	existing.SavedCardHolderName = coalesceString(user.SavedCardHolderName, existing.SavedCardHolderName)
	existing.SavedCardNumber = coalesceString(user.SavedCardNumber, existing.SavedCardNumber)
	existing.SavedCardExpiry = coalesceString(user.SavedCardExpiry, existing.SavedCardExpiry)
	existing.BankName = coalesceString(user.BankName, existing.BankName)
	existing.BankAccountNumber = coalesceString(user.BankAccountNumber, existing.BankAccountNumber)
	existing.BankBookImageName = coalesceString(user.BankBookImageName, existing.BankBookImageName)
	existing.BankBookImageURL = coalesceString(user.BankBookImageURL, existing.BankBookImageURL)
	existing.IdentityCardImageName = coalesceString(user.IdentityCardImageName, existing.IdentityCardImageName)
	existing.IdentityCardImageURL = coalesceString(user.IdentityCardImageURL, existing.IdentityCardImageURL)
	existing.NotifyOrders = user.NotifyOrders
	existing.NotifyMarketplace = user.NotifyMarketplace

	if existing.Role == "seller" && strings.TrimSpace(existing.Slug) == "" {
		slug, slugErr := r.uniqueSlug(coalesceString(existing.StoreName, existing.Name), existing.ID)
		if slugErr != nil {
			return nil, slugErr
		}
		existing.Slug = slug
	}

	if err := r.db.Save(&existing).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (r *UserRepository) dbSaveAndReturn(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) uniqueSlug(base string, excludeUserID string) (string, error) {
	candidate := slugify(base)
	if candidate == "" {
		candidate = fmt.Sprintf("seller-%d", time.Now().Unix())
	}

	currentCandidate := candidate
	for suffix := 1; ; suffix++ {
		var count int64
		query := r.db.Model(&models.User{}).Where("slug = ?", currentCandidate)
		if strings.TrimSpace(excludeUserID) != "" {
			query = query.Where("id <> ?", strings.TrimSpace(excludeUserID))
		}

		if err := query.Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return currentCandidate, nil
		}

		currentCandidate = fmt.Sprintf("%s-%d", candidate, suffix+1)
	}
}

func coalesceString(next string, fallback string) string {
	if strings.TrimSpace(next) == "" {
		return fallback
	}

	return next
}

func deriveNameFromEmail(email string, role string) string {
	prefix := strings.TrimSpace(strings.Split(strings.TrimSpace(email), "@")[0])
	if prefix == "" {
		if role == "seller" {
			return "CodeBazaar Seller"
		}
		return "ผู้ใช้ CodeBazaar"
	}

	parts := strings.Fields(strings.NewReplacer(".", " ", "_", " ", "-", " ").Replace(prefix))
	for index, part := range parts {
		if part == "" {
			continue
		}
		parts[index] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	return strings.Join(parts, " ")
}

func resolveProvider(provider string, userID string) string {
	normalizedProvider := strings.TrimSpace(strings.ToLower(provider))
	if normalizedProvider != "" {
		return normalizedProvider
	}

	switch {
	case strings.HasPrefix(strings.TrimSpace(userID), "github-"):
		return "github"
	case strings.HasPrefix(strings.TrimSpace(userID), "google-"):
		return "google"
	default:
		return "google"
	}
}

func resolveRole(role string, provider string) string {
	normalizedRole := strings.TrimSpace(strings.ToLower(role))
	if normalizedRole != "" {
		return normalizedRole
	}

	if provider == "github" {
		return "seller"
	}

	return "buyer"
}

func mergeRole(existing string, next string) string {
	normalizedExisting := strings.TrimSpace(strings.ToLower(existing))
	normalizedNext := strings.TrimSpace(strings.ToLower(next))

	switch {
	case normalizedNext == "":
		return existing
	case normalizedExisting == "":
		return normalizedNext
	case normalizedExisting == "seller" && normalizedNext == "buyer":
		return normalizedExisting
	default:
		return normalizedNext
	}
}

func mergeProvider(existing string, next string) string {
	normalizedExisting := strings.TrimSpace(strings.ToLower(existing))
	normalizedNext := strings.TrimSpace(strings.ToLower(next))

	switch {
	case normalizedNext == "":
		return existing
	case normalizedExisting == "":
		return normalizedNext
	case normalizedNext == "guest":
		return normalizedExisting
	default:
		return normalizedNext
	}
}

func slugify(value string) string {
	normalized := strings.TrimSpace(strings.ToLower(value))
	normalized = strings.NewReplacer(".", "-", "_", "-", " ", "-").Replace(normalized)
	normalized = strings.Join(strings.Fields(normalized), "-")
	normalized = strings.Trim(normalized, "-")
	for strings.Contains(normalized, "--") {
		normalized = strings.ReplaceAll(normalized, "--", "-")
	}

	if normalized == "" {
		return fmt.Sprintf("seller-%d", time.Now().Unix())
	}

	return normalized
}

func createRepositoryID(prefix string) string {
	buffer := make([]byte, 4)
	if _, err := rand.Read(buffer); err != nil {
		return prefix + "-" + time.Now().Format("20060102150405.000000000")
	}

	return prefix + "-" + time.Now().Format("20060102150405.000000000") + "-" + hex.EncodeToString(buffer)
}
