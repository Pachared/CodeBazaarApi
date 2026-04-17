package repositories

import (
	"github.com/Pachared/CodeBazaarApi/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) ListPublishedWithSeller() ([]models.Product, error) {
	var products []models.Product
	err := r.db.
		Preload("Seller").
		Where("status = ?", "publish").
		Find(&products).Error

	return products, err
}

func (r *ProductRepository) GetPublishedByID(productID string) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Seller").
		First(&product, "id = ? AND status = ?", productID, "publish").Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) ListByIDs(productIDs []string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.
		Preload("Seller").
		Where("id IN ?", productIDs).
		Find(&products).Error

	return products, err
}

func (r *ProductRepository) ListPublishedBySellerSlug(slug string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.
		Preload("Seller").
		Joins("JOIN users ON users.id = products.seller_id").
		Where("products.status = ? AND users.slug = ?", "publish", slug).
		Find(&products).Error

	return products, err
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}
