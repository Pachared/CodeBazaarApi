package repositories

import (
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type SellerOrderRow struct {
	ID              string
	OrderID         string
	ProductID       string
	ProductTitle    string
	ProductCategory string
	BuyerName       string
	BuyerEmail      string
	PurchasedAt     time.Time
	Amount          int64
	PaymentMethod   string
	LicenseLabel    string
	DeliveryLabel   string
	Status          string
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepository) CreateDownload(tx *gorm.DB, item *models.DownloadItem) error {
	return tx.Create(item).Error
}

func (r *OrderRepository) ListSellerOrders(sellerID string) ([]SellerOrderRow, error) {
	var rows []SellerOrderRow

	query := r.db.Model(&models.OrderItem{}).
		Select(`
			order_items.id AS id,
			orders.order_number AS order_id,
			order_items.product_id AS product_id,
			order_items.title AS product_title,
			order_items.category AS product_category,
			orders.customer_name AS buyer_name,
			orders.customer_email AS buyer_email,
			orders.purchased_at AS purchased_at,
			order_items.price AS amount,
			orders.payment_method AS payment_method,
			order_items.license AS license_label,
			order_items.delivery_label AS delivery_label,
			orders.status AS status
		`).
		Joins("JOIN orders ON orders.id = order_items.order_id")

	if sellerID != "" {
		query = query.Where("order_items.seller_id = ?", sellerID)
	}

	err := query.
		Order("orders.purchased_at DESC").
		Scan(&rows).Error

	return rows, err
}

func (r *OrderRepository) ListDownloads(userID string) ([]models.DownloadItem, error) {
	var items []models.DownloadItem
	err := r.db.
		Where("user_id = ?", userID).
		Order("purchased_at DESC").
		Find(&items).Error

	return items, err
}

func (r *OrderRepository) FindDownloadForUser(userID string, itemID string) (*models.DownloadItem, error) {
	var item models.DownloadItem
	err := r.db.First(&item, "id = ? AND user_id = ?", itemID, userID).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *OrderRepository) SaveDownload(item *models.DownloadItem) error {
	return r.db.Save(item).Error
}
