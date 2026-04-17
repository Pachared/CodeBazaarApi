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

const (
	defaultBuyerID     = "usr_buyer_demo"
	defaultBuyerEmail  = "buyer.demo@codebazaar.local"
	defaultSellerID    = "usr_seller_demo"
	defaultSellerEmail = "seller.demo@codebazaar.local"
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

func (r *UserRepository) FindOrCreateDemoBuyer(intent string) (*models.User, error) {
	user := &models.User{
		ID:                  defaultBuyerID,
		Name:                "ผู้ซื้อทดลอง",
		Email:               defaultBuyerEmail,
		Role:                "buyer",
		Provider:            "google",
		IsMock:              true,
		PhoneNumber:         "0812345678",
		SavedCardHolderName: "Pachara Demo",
		SavedCardNumber:     "4111 1111 1111 1111",
		SavedCardExpiry:     "12/28",
		NotifyOrders:        true,
		NotifyMarketplace:   true,
	}

	if intent == "register" {
		user.Name = "สมาชิกทดลอง CodeBazaar"
	}

	return r.firstOrCreateByEmail(user)
}

func (r *UserRepository) FindOrCreateDemoSeller() (*models.User, error) {
	user := &models.User{
		ID:                    defaultSellerID,
		Slug:                  "codebazaar-seller-demo",
		Name:                  "CodeBazaar Seller Demo",
		Email:                 defaultSellerEmail,
		Role:                  "seller",
		Provider:              "github",
		IsMock:                true,
		PhoneNumber:           "0898765432",
		StoreName:             "CodeBazaar Seller Demo",
		BankName:              "ธนาคารกสิกรไทย",
		BankAccountNumber:     "123-4-56789-0",
		BankBookImageName:     "bank-book-demo.png",
		BankBookImageURL:      "https://example.com/bank-book-demo.png",
		IdentityCardImageName: "identity-card-demo.png",
		IdentityCardImageURL:  "https://example.com/identity-card-demo.png",
		NotifyOrders:          true,
		NotifyMarketplace:     true,
	}

	return r.firstOrCreateByEmail(user)
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
		existingUser.ID = coalesceString(normalizedID, existingUser.ID)
		existingUser.Name = coalesceString(strings.TrimSpace(name), existingUser.Name)
		existingUser.Email = coalesceString(normalizedEmail, existingUser.Email)
		existingUser.Provider = coalesceString(normalizedProvider, existingUser.Provider)
		existingUser.Role = coalesceString(normalizedRole, existingUser.Role)
		existingUser.IsMock = false

		if existingUser.Role == "seller" {
			existingUser.Slug = coalesceString(existingUser.Slug, slugify(existingUser.Name))
			existingUser.StoreName = coalesceString(existingUser.StoreName, existingUser.Name)
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
		IsMock:            false,
		NotifyOrders:      true,
		NotifyMarketplace: true,
	}

	if user.Role == "seller" {
		user.Slug = slugify(displayName)
		user.StoreName = displayName
	}

	return r.firstOrCreateByEmail(user)
}

func (r *UserRepository) FindOrCreateBuyerByEmail(name string, email string, phone string) (*models.User, error) {
	normalizedEmail := strings.TrimSpace(strings.ToLower(email))
	if normalizedEmail == "" {
		return r.FindOrCreateDemoBuyer("login")
	}

	user := &models.User{
		ID:                createRepositoryID("usr_buyer"),
		Name:              strings.TrimSpace(name),
		Email:             normalizedEmail,
		Role:              "buyer",
		Provider:          "google",
		IsMock:            true,
		PhoneNumber:       strings.TrimSpace(phone),
		NotifyOrders:      true,
		NotifyMarketplace: true,
	}

	if user.Name == "" {
		user.Name = "ผู้ซื้อ CodeBazaar"
	}

	return r.firstOrCreateByEmail(user)
}

func (r *UserRepository) ResolveOrDefaultBuyer(currentUser *models.User) (*models.User, error) {
	if currentUser != nil {
		return currentUser, nil
	}

	return r.FindOrCreateDemoBuyer("login")
}

func (r *UserRepository) ResolveOrDefaultSeller(currentUser *models.User) (*models.User, error) {
	if currentUser != nil && currentUser.Role == "seller" {
		return currentUser, nil
	}

	return r.FindOrCreateDemoSeller()
}

func (r *UserRepository) firstOrCreateByEmail(seed *models.User) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", strings.TrimSpace(strings.ToLower(seed.Email))).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.Create(seed).Error; err != nil {
				return nil, err
			}
			return seed, nil
		}

		return nil, err
	}

	user.Name = coalesceString(seed.Name, user.Name)
	user.Role = coalesceString(seed.Role, user.Role)
	user.Provider = coalesceString(seed.Provider, user.Provider)
	user.IsMock = seed.IsMock
	user.Slug = coalesceString(seed.Slug, user.Slug)
	user.PhoneNumber = coalesceString(seed.PhoneNumber, user.PhoneNumber)
	user.StoreName = coalesceString(seed.StoreName, user.StoreName)
	user.SavedCardHolderName = coalesceString(seed.SavedCardHolderName, user.SavedCardHolderName)
	user.SavedCardNumber = coalesceString(seed.SavedCardNumber, user.SavedCardNumber)
	user.SavedCardExpiry = coalesceString(seed.SavedCardExpiry, user.SavedCardExpiry)
	user.BankName = coalesceString(seed.BankName, user.BankName)
	user.BankAccountNumber = coalesceString(seed.BankAccountNumber, user.BankAccountNumber)
	user.BankBookImageName = coalesceString(seed.BankBookImageName, user.BankBookImageName)
	user.BankBookImageURL = coalesceString(seed.BankBookImageURL, user.BankBookImageURL)
	user.IdentityCardImageName = coalesceString(seed.IdentityCardImageName, user.IdentityCardImageName)
	user.IdentityCardImageURL = coalesceString(seed.IdentityCardImageURL, user.IdentityCardImageURL)
	user.NotifyOrders = seed.NotifyOrders
	user.NotifyMarketplace = seed.NotifyMarketplace

	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
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
