package services

import (
	"sort"
	"strings"
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/models"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
)

type SellerService struct {
	userRepository    *repositories.UserRepository
	productRepository *repositories.ProductRepository
	orderRepository   *repositories.OrderRepository
}

func NewSellerService(
	userRepository *repositories.UserRepository,
	productRepository *repositories.ProductRepository,
	orderRepository *repositories.OrderRepository,
) *SellerService {
	return &SellerService{
		userRepository:    userRepository,
		productRepository: productRepository,
		orderRepository:   orderRepository,
	}
}

func (s *SellerService) OpenSellerAccount() (*contracts.AuthActionResponse, error) {
	user, err := s.userRepository.FindOrCreateDemoSeller()
	if err != nil {
		return nil, err
	}

	return &contracts.AuthActionResponse{
		Title:       "เปิดบัญชีผู้ขายสำเร็จ",
		Description: "เข้าสู่ระบบด้วยบัญชีผู้ขายทดลองเรียบร้อยแล้ว และพร้อมใช้งาน Seller Studio ต่อได้ทันที",
		Session:     toAuthSessionUser(user),
	}, nil
}

func (s *SellerService) SubmitListing(currentUser *models.User, input contracts.SellerListingRequest) (*contracts.SellerListingResponse, error) {
	seller, err := s.userRepository.ResolveOrDefaultSeller(currentUser)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		ID:                    "listing-" + strings.ReplaceAll(createStableID("product"), "_", "-"),
		SellerID:              seller.ID,
		AssetType:             input.AssetType,
		Status:                input.Mode,
		CategoryID:            input.CategoryID,
		Category:              categoryLabel(input.CategoryID),
		Title:                 strings.TrimSpace(input.Title),
		Summary:               strings.TrimSpace(input.Summary),
		FullDescription:       strings.TrimSpace(input.Description),
		Price:                 input.Price,
		Rating:                0,
		Sales:                 0,
		Tags:                  dedupeStrings(append([]string{categoryLabel(input.CategoryID)}, input.Stack...), 0),
		Stack:                 dedupeStrings(input.Stack, 0),
		FeatureHighlights:     dedupeStrings(input.Highlights, 0),
		IncludedItems:         dedupeStrings(input.IncludedFiles, 0),
		IdealFor:              dedupeStrings(input.IdealFor, 0),
		SupportInfo:           strings.TrimSpace(input.SupportInfo),
		VersionLabel:          strings.TrimSpace(input.Version),
		FileFormatLabel:       buildFileFormatLabel(input.DocumentationIncluded),
		UpdatedDaysAgo:        0,
		Delivery:              deliveryLabel(input.InstantDelivery),
		License:               licenseLabel(input.LicenseID),
		LicenseID:             input.LicenseID,
		Verified:              seller.IdentityCardImageName != "" && seller.BankBookImageName != "",
		DemoURL:               strings.TrimSpace(input.DemoURL),
		SupportURL:            strings.TrimSpace(input.SupportURL),
		PackageFileName:       strings.TrimSpace(input.PackageFileName),
		CoverFileName:         strings.TrimSpace(input.CoverFileName),
		DocsFileName:          strings.TrimSpace(input.DocsFileName),
		InstantDelivery:       input.InstantDelivery,
		SourceIncluded:        input.SourceIncluded,
		DocumentationIncluded: input.DocumentationIncluded,
		TimestampedModel: models.TimestampedModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if product.VersionLabel == "" {
		product.VersionLabel = "v1.0"
	}

	if err := s.productRepository.Create(product); err != nil {
		return nil, err
	}

	title := "บันทึกร่างรายการแล้ว"
	description := product.Title + " ถูกบันทึกเป็นฉบับร่างเรียบร้อยแล้ว"
	if input.Mode == "publish" {
		title = "ส่งขึ้นรายการขายแล้ว"
		description = product.Title + " ถูกส่งขึ้นพื้นที่ขายเรียบร้อยแล้ว และพร้อมต่อ workflow อนุมัติรายการจริง"
	}

	return &contracts.SellerListingResponse{
		Title:       title,
		Description: description,
		ListingID:   product.ID,
		Status:      input.Mode,
	}, nil
}

func (s *SellerService) ListSellerOrders(currentUser *models.User) ([]contracts.SellerOrderResponse, error) {
	sellerID := ""
	if currentUser != nil && currentUser.Role == "seller" {
		sellerID = currentUser.ID
	}

	rows, err := s.orderRepository.ListSellerOrders(sellerID)
	if err != nil {
		return nil, err
	}

	responses := make([]contracts.SellerOrderResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, contracts.SellerOrderResponse{
			ID:                 row.ID,
			OrderID:            row.OrderID,
			ProductID:          row.ProductID,
			ProductTitle:       row.ProductTitle,
			ProductCategory:    row.ProductCategory,
			BuyerName:          row.BuyerName,
			BuyerEmail:         row.BuyerEmail,
			PurchasedAt:        row.PurchasedAt.Format(time.RFC3339),
			Amount:             row.Amount,
			PaymentMethodLabel: paymentMethodLabel(row.PaymentMethod),
			LicenseLabel:       row.LicenseLabel,
			DeliveryLabel:      row.DeliveryLabel,
			StatusLabel:        orderStatusLabel(row.Status),
		})
	}

	sort.SliceStable(responses, func(left int, right int) bool {
		return responses[left].PurchasedAt > responses[right].PurchasedAt
	})

	return responses, nil
}
