package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/models"
)

var (
	categoryLabelMap = map[string]string{
		"marketplace":   "มาร์เก็ตเพลส",
		"dashboard":     "แดชบอร์ด",
		"landing-page":  "หน้าเปิดตัว",
		"saas":          "SaaS เริ่มต้น",
		"design-system": "ระบบดีไซน์",
	}
	licenseLabelMap = map[string]string{
		"personal":   "ใช้งานส่วนตัว",
		"commercial": "ใช้งานเชิงพาณิชย์",
		"resale":     "ขายต่อได้",
	}
	paymentMethodLabelMap = map[string]string{
		"promptpay":     "พร้อมเพย์ QR",
		"card":          "บัตรเครดิต / เดบิต",
		"bank-transfer": "โอนผ่านบัญชีธนาคาร",
	}
	alphaNumOnlyPattern = regexp.MustCompile(`[^a-z0-9]+`)
)

func createStableID(prefix string) string {
	buffer := make([]byte, 4)
	if _, err := rand.Read(buffer); err != nil {
		return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	}

	return fmt.Sprintf("%s-%d-%s", prefix, time.Now().UnixNano(), hex.EncodeToString(buffer))
}

func createOrderNumber() string {
	return fmt.Sprintf("CB-%08d", time.Now().UnixNano()%100000000)
}

func categoryLabel(categoryID string) string {
	if label, ok := categoryLabelMap[categoryID]; ok {
		return label
	}

	return categoryID
}

func licenseLabel(licenseID string) string {
	if label, ok := licenseLabelMap[licenseID]; ok {
		return label
	}

	return licenseID
}

func paymentMethodLabel(method string) string {
	if label, ok := paymentMethodLabelMap[method]; ok {
		return label
	}

	return method
}

func orderStatusLabel(status string) string {
	switch status {
	case "paid":
		return "ชำระเงินสำเร็จ"
	case "pending":
		return "รอตรวจสอบสลิป"
	default:
		return status
	}
}

func updatedLabel(days int) string {
	switch {
	case days <= 0:
		return "อัปเดตวันนี้"
	case days == 1:
		return "อัปเดต 1 วันที่แล้ว"
	default:
		return fmt.Sprintf("อัปเดต %d วันที่แล้ว", days)
	}
}

func buildFileFormatLabel(documentationIncluded bool) string {
	if documentationIncluded {
		return "ZIP + เอกสารประกอบ"
	}

	return "ZIP"
}

func deliveryLabel(instant bool) string {
	if instant {
		return "ดาวน์โหลดได้ทันที"
	}

	return "รออนุมัติจากผู้ขาย"
}

func buildDownloadFileName(title string, productID string) string {
	normalized := strings.ToLower(strings.TrimSpace(title))
	normalized = alphaNumOnlyPattern.ReplaceAllString(normalized, "-")
	normalized = strings.Trim(normalized, "-")
	if normalized == "" {
		normalized = "codebazaar-" + strings.TrimSpace(productID)
	}

	return normalized + "-package.zip"
}

func buildFileSizeLabel(price int64) string {
	megabytes := price / 170
	if megabytes < 6 {
		megabytes = 6
	}

	return fmt.Sprintf("%d MB", megabytes)
}

func dedupeStrings(values []string, limit int) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.TrimSpace(value)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		unique = append(unique, normalized)
		if limit > 0 && len(unique) >= limit {
			break
		}
	}

	return unique
}

func buildSellerSummary(products []models.Product) string {
	categories := dedupeStrings(extractCategories(products), 0)

	switch len(categories) {
	case 0:
		return "รวมรายการซอร์สโค้ดและเทมเพลตที่พร้อมนำไปต่อยอดในระบบจริง"
	case 1:
		return fmt.Sprintf("รวมผลงานในหมวด %s ที่พร้อมให้ดูรายละเอียดและกดซื้อได้จากหน้าร้านของผู้ขายรายนี้", categories[0])
	default:
		return fmt.Sprintf("รวมผลงานในหมวด %s และ %s พร้อมรายการซอร์สโค้ดและเทมเพลตที่นำไปต่อยอดได้ทันที", categories[0], categories[1])
	}
}

func extractCategories(products []models.Product) []string {
	categories := make([]string, 0, len(products))
	for _, product := range products {
		categories = append(categories, product.Category)
	}

	return categories
}

func toAuthSessionUser(user *models.User) *contracts.AuthSessionUser {
	if user == nil {
		return nil
	}

	return &contracts.AuthSessionUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Provider: user.Provider,
		AuthProfileFields: contracts.AuthProfileFields{
			PhoneNumber:           user.PhoneNumber,
			StoreName:             user.StoreName,
			SavedCardHolderName:   user.SavedCardHolderName,
			SavedCardNumber:       user.SavedCardNumber,
			SavedCardExpiry:       user.SavedCardExpiry,
			BankName:              user.BankName,
			BankAccountNumber:     user.BankAccountNumber,
			BankBookImageName:     user.BankBookImageName,
			BankBookImageURL:      user.BankBookImageURL,
			IdentityCardImageName: user.IdentityCardImageName,
			IdentityCardImageURL:  user.IdentityCardImageURL,
			NotifyOrders:          user.NotifyOrders,
			NotifyMarketplace:     user.NotifyMarketplace,
		},
	}
}

func requireCurrentUser(currentUser *models.User) (*models.User, error) {
	if currentUser == nil || strings.TrimSpace(currentUser.ID) == "" {
		return nil, httpx.NewAppError(http.StatusUnauthorized, "กรุณาเข้าสู่ระบบก่อนใช้งานส่วนนี้")
	}

	return currentUser, nil
}

func requireSellerUser(currentUser *models.User) (*models.User, error) {
	user, err := requireCurrentUser(currentUser)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(strings.ToLower(user.Role)) != "seller" {
		return nil, httpx.NewAppError(http.StatusForbidden, "บัญชีนี้ยังไม่ได้เปิดใช้งานสิทธิ์ผู้ขาย")
	}

	return user, nil
}

func toProductResponse(product models.Product) contracts.ProductResponse {
	return contracts.ProductResponse{
		ID:                product.ID,
		CategoryID:        product.CategoryID,
		Title:             product.Title,
		Summary:           product.Summary,
		FullDescription:   product.FullDescription,
		Category:          product.Category,
		Price:             product.Price,
		Rating:            product.Rating,
		Sales:             product.Sales,
		Tags:              append([]string{}, product.Tags...),
		Stack:             append([]string{}, product.Stack...),
		FeatureHighlights: append([]string{}, product.FeatureHighlights...),
		IncludedItems:     append([]string{}, product.IncludedItems...),
		IdealFor:          append([]string{}, product.IdealFor...),
		SupportInfo:       product.SupportInfo,
		VersionLabel:      product.VersionLabel,
		FileFormatLabel:   product.FileFormatLabel,
		AuthorName:        product.Seller.Name,
		AuthorSlug:        product.Seller.Slug,
		UpdatedAt:         updatedLabel(product.UpdatedDaysAgo),
		UpdatedDaysAgo:    product.UpdatedDaysAgo,
		Delivery:          product.Delivery,
		License:           product.License,
		LicenseID:         product.LicenseID,
		Verified:          product.Verified,
	}
}

func toMarketplaceSellerResponses(products []models.Product) []contracts.MarketplaceSellerResponse {
	groupedProducts := make(map[string][]models.Product)
	for _, product := range products {
		groupedProducts[product.Seller.Slug] = append(groupedProducts[product.Seller.Slug], product)
	}

	sellers := make([]contracts.MarketplaceSellerResponse, 0, len(groupedProducts))
	for slug, sellerProducts := range groupedProducts {
		if len(sellerProducts) == 0 {
			continue
		}

		firstProduct := sellerProducts[0]
		totalSales := 0
		startingPrice := sellerProducts[0].Price
		verifiedCount := 0
		latestUpdateDaysAgo := sellerProducts[0].UpdatedDaysAgo
		stacks := make([]string, 0)
		categories := make([]string, 0)

		for _, product := range sellerProducts {
			totalSales += product.Sales
			if product.Price < startingPrice {
				startingPrice = product.Price
			}
			if product.Verified {
				verifiedCount++
			}
			if product.UpdatedDaysAgo < latestUpdateDaysAgo {
				latestUpdateDaysAgo = product.UpdatedDaysAgo
			}
			stacks = append(stacks, product.Stack...)
			categories = append(categories, product.Category)
		}

		sellers = append(sellers, contracts.MarketplaceSellerResponse{
			Slug:                slug,
			Name:                firstProduct.Seller.Name,
			Summary:             buildSellerSummary(sellerProducts),
			ProductCount:        len(sellerProducts),
			TotalSales:          totalSales,
			StartingPrice:       startingPrice,
			VerifiedCount:       verifiedCount,
			Categories:          dedupeStrings(categories, 0),
			Stacks:              dedupeStrings(stacks, 5),
			LatestUpdateDaysAgo: latestUpdateDaysAgo,
		})
	}

	sort.SliceStable(sellers, func(left int, right int) bool {
		if sellers[left].VerifiedCount != sellers[right].VerifiedCount {
			return sellers[left].VerifiedCount > sellers[right].VerifiedCount
		}
		return sellers[left].TotalSales > sellers[right].TotalSales
	})

	return sellers
}

func toDownloadItemResponse(item models.DownloadItem) contracts.DownloadLibraryItemResponse {
	var lastDownloadedAt *string
	if item.LastDownloadedAt != nil {
		value := item.LastDownloadedAt.Format(time.RFC3339)
		lastDownloadedAt = &value
	}

	return contracts.DownloadLibraryItemResponse{
		LibraryItemID:    item.ID,
		OrderID:          item.OrderID,
		ID:               item.ProductID,
		Title:            item.Title,
		Category:         item.Category,
		AuthorName:       item.AuthorName,
		Price:            item.Price,
		License:          item.License,
		PurchasedAt:      item.PurchasedAt.Format(time.RFC3339),
		PaymentMethod:    item.PaymentMethod,
		Status:           item.Status,
		VersionLabel:     item.VersionLabel,
		FileName:         item.FileName,
		FileSizeLabel:    item.FileSizeLabel,
		DownloadsCount:   item.DownloadsCount,
		LastDownloadedAt: lastDownloadedAt,
	}
}

func toCookieConsentResponse(consent *models.CookieConsent) *contracts.CookieConsentResponse {
	if consent == nil {
		return nil
	}

	return &contracts.CookieConsentResponse{
		Status: consent.Status,
		Preferences: contracts.CookiePreferences{
			Necessary:   consent.Preferences.Necessary,
			Preferences: consent.Preferences.Preferences,
			Analytics:   consent.Preferences.Analytics,
			Marketing:   consent.Preferences.Marketing,
		},
		UpdatedAt: consent.ConsentUpdatedAt.Format(time.RFC3339),
	}
}
