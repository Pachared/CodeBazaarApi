package services

import (
	"net/http"
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/models"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
	"gorm.io/gorm"
)

type CheckoutService struct {
	db                *gorm.DB
	userRepository    *repositories.UserRepository
	productRepository *repositories.ProductRepository
	orderRepository   *repositories.OrderRepository
}

func NewCheckoutService(
	db *gorm.DB,
	userRepository *repositories.UserRepository,
	productRepository *repositories.ProductRepository,
	orderRepository *repositories.OrderRepository,
) *CheckoutService {
	return &CheckoutService{
		db:                db,
		userRepository:    userRepository,
		productRepository: productRepository,
		orderRepository:   orderRepository,
	}
}

func (s *CheckoutService) SubmitOrder(currentUser *models.User, input contracts.CheckoutSubmitInput) (*contracts.CheckoutSubmitResponse, error) {
	if len(input.Items) == 0 {
		return nil, httpx.NewAppError(http.StatusBadRequest, "ยังไม่มีรายการสินค้าสำหรับชำระเงิน")
	}

	productIDs := make([]string, 0, len(input.Items))
	for _, item := range input.Items {
		productIDs = append(productIDs, item.ID)
	}

	products, err := s.productRepository.ListByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	productByID := make(map[string]models.Product, len(products))
	for _, product := range products {
		productByID[product.ID] = product
	}

	for _, item := range input.Items {
		if _, ok := productByID[item.ID]; !ok {
			return nil, httpx.NewAppError(http.StatusBadRequest, "มีรายการสินค้าบางชิ้นไม่พบในระบบแล้ว กรุณารีเฟรชข้อมูลก่อนลองใหม่")
		}
	}

	buyer, err := s.userRepository.FindOrCreateBuyerByEmail(input.CustomerName, input.CustomerEmail, input.CustomerPhone)
	if err != nil {
		return nil, err
	}

	if currentUser != nil && currentUser.ID != "" {
		buyer = currentUser
	}

	order := &models.Order{
		ID:                     createStableID("ord"),
		OrderNumber:            createOrderNumber(),
		BuyerID:                buyer.ID,
		CustomerName:           input.CustomerName,
		CustomerEmail:          input.CustomerEmail,
		CustomerPhone:          input.CustomerPhone,
		CompanyName:            input.CompanyName,
		TaxID:                  input.TaxID,
		Note:                   input.Note,
		PaymentMethod:          input.PaymentMethod,
		ReceivePurchaseUpdates: input.ReceivePurchaseUpdates,
		RequireInvoice:         input.RequireInvoice,
		Subtotal:               input.Subtotal,
		Total:                  input.Total,
		Status:                 "paid",
		PaymentDetails: models.CheckoutPaymentDetails{
			PromptPay:    toModelPromptPay(input.PaymentDetails.PromptPay),
			Card:         toModelCard(input.PaymentDetails.Card),
			BankTransfer: toModelBankTransfer(input.PaymentDetails.BankTransfer),
		},
		PurchasedAt: time.Now(),
	}

	if buyer.PhoneNumber != input.CustomerPhone {
		buyer.PhoneNumber = input.CustomerPhone
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.orderRepository.CreateOrder(tx, order); err != nil {
			return err
		}

		for _, item := range input.Items {
			product := productByID[item.ID]

			orderItem := &models.OrderItem{
				ID:            createStableID("item"),
				OrderID:       order.ID,
				ProductID:     product.ID,
				SellerID:      product.SellerID,
				Title:         product.Title,
				Category:      product.Category,
				AuthorName:    product.Seller.Name,
				Price:         product.Price,
				License:       product.License,
				DeliveryLabel: product.Delivery,
			}

			if err := tx.Create(orderItem).Error; err != nil {
				return err
			}

			downloadItem := &models.DownloadItem{
				ID:             createStableID("dl"),
				UserID:         buyer.ID,
				OrderID:        order.OrderNumber,
				OrderItemID:    orderItem.ID,
				ProductID:      product.ID,
				Title:          product.Title,
				Category:       product.Category,
				AuthorName:     product.Seller.Name,
				Price:          product.Price,
				License:        product.License,
				PaymentMethod:  input.PaymentMethod,
				PurchasedAt:    order.PurchasedAt,
				Status:         "ready",
				VersionLabel:   product.VersionLabel,
				FileName:       buildDownloadFileName(product.Title, product.ID),
				FileSizeLabel:  buildFileSizeLabel(product.Price),
				DownloadsCount: 0,
			}

			if err := s.orderRepository.CreateDownload(tx, downloadItem); err != nil {
				return err
			}
		}

		return tx.Save(buyer).Error
	})
	if err != nil {
		return nil, err
	}

	return &contracts.CheckoutSubmitResponse{
		Title:         "ชำระเงินสำเร็จ",
		Description:   "รายการของคุณถูกยืนยันแล้ว และพร้อมดาวน์โหลดได้ทันทีในระบบจริง",
		OrderID:       order.OrderNumber,
		PaymentMethod: input.PaymentMethod,
		Total:         input.Total,
		Status:        "paid",
	}, nil
}

func toModelPromptPay(source *contracts.PromptPayPaymentDetails) *models.PromptPayPaymentDetails {
	if source == nil {
		return nil
	}

	return &models.PromptPayPaymentDetails{
		AccountName:   source.AccountName,
		PromptPayID:   source.PromptPayID,
		ReferenceCode: source.ReferenceCode,
	}
}

func toModelCard(source *contracts.CardPaymentDetails) *models.CardPaymentDetails {
	if source == nil {
		return nil
	}

	return &models.CardPaymentDetails{
		CardHolderName:   source.CardHolderName,
		CardNumberMasked: source.CardNumberMasked,
		Expiry:           source.Expiry,
	}
}

func toModelBankTransfer(source *contracts.BankTransferPaymentDetails) *models.BankTransferPaymentDetails {
	if source == nil {
		return nil
	}

	return &models.BankTransferPaymentDetails{
		BankName:          source.BankName,
		AccountName:       source.AccountName,
		AccountNumber:     source.AccountNumber,
		TransferReference: source.TransferReference,
		SlipImageName:     source.SlipImageName,
	}
}
