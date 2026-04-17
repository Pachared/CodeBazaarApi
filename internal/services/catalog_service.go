package services

import (
	"sort"
	"strings"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/models"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
)

type ProductFilter struct {
	Query        string
	Category     string
	License      string
	Price        string
	Sort         string
	VerifiedOnly bool
	Stacks       []string
}

type CatalogService struct {
	productRepository *repositories.ProductRepository
}

func NewCatalogService(productRepository *repositories.ProductRepository) *CatalogService {
	return &CatalogService{productRepository: productRepository}
}

func (s *CatalogService) ListFeaturedProducts() ([]contracts.ProductResponse, error) {
	products, err := s.productRepository.ListPublishedWithSeller()
	if err != nil {
		return nil, err
	}

	sortProducts(products, "featured")

	responses := make([]contracts.ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, toProductResponse(product))
	}

	return responses, nil
}

func (s *CatalogService) GetProductByID(productID string) (*contracts.ProductResponse, error) {
	product, err := s.productRepository.GetPublishedByID(productID)
	if err != nil {
		return nil, httpx.NewAppError(404, "ไม่พบรายการสินค้าที่คุณต้องการดูรายละเอียด")
	}

	response := toProductResponse(*product)
	return &response, nil
}

func (s *CatalogService) ListProducts(filter ProductFilter) ([]contracts.ProductResponse, error) {
	products, err := s.productRepository.ListPublishedWithSeller()
	if err != nil {
		return nil, err
	}

	filteredProducts := make([]models.Product, 0, len(products))
	for _, product := range products {
		if !matchesFilter(product, filter) {
			continue
		}
		filteredProducts = append(filteredProducts, product)
	}

	sortProducts(filteredProducts, filter.Sort)

	responses := make([]contracts.ProductResponse, 0, len(filteredProducts))
	for _, product := range filteredProducts {
		responses = append(responses, toProductResponse(product))
	}

	return responses, nil
}

func (s *CatalogService) ListSellers() ([]contracts.MarketplaceSellerResponse, error) {
	products, err := s.productRepository.ListPublishedWithSeller()
	if err != nil {
		return nil, err
	}

	return toMarketplaceSellerResponses(products), nil
}

func (s *CatalogService) GetSellerBySlug(slug string) (*contracts.MarketplaceSellerResponse, error) {
	sellers, err := s.ListSellers()
	if err != nil {
		return nil, err
	}

	for _, seller := range sellers {
		if seller.Slug == slug {
			return &seller, nil
		}
	}

	return nil, httpx.NewAppError(404, "ไม่พบผู้ขายรายนี้ในระบบ")
}

func (s *CatalogService) ListProductsBySellerSlug(slug string) ([]contracts.ProductResponse, error) {
	products, err := s.productRepository.ListPublishedBySellerSlug(slug)
	if err != nil {
		return nil, err
	}

	sortProducts(products, "featured")

	responses := make([]contracts.ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, toProductResponse(product))
	}

	return responses, nil
}

func matchesFilter(product models.Product, filter ProductFilter) bool {
	normalizedQuery := strings.TrimSpace(strings.ToLower(filter.Query))
	if normalizedQuery != "" {
		searchSpace := strings.ToLower(strings.Join(append([]string{
			product.Title,
			product.Summary,
			product.Category,
			product.Seller.Name,
		}, product.Stack...), " "))

		if !strings.Contains(searchSpace, normalizedQuery) {
			return false
		}
	}

	if filter.Category != "" && filter.Category != "all" && product.CategoryID != filter.Category {
		return false
	}

	if filter.License != "" && filter.License != "all" && product.LicenseID != filter.License {
		return false
	}

	if filter.VerifiedOnly && !product.Verified {
		return false
	}

	if len(filter.Stacks) > 0 {
		productStacks := make(map[string]struct{}, len(product.Stack))
		for _, stack := range product.Stack {
			productStacks[stack] = struct{}{}
		}

		for _, stack := range filter.Stacks {
			if _, ok := productStacks[stack]; !ok {
				return false
			}
		}
	}

	switch filter.Price {
	case "under-1500":
		return product.Price < 1500
	case "1500-2500":
		return product.Price >= 1500 && product.Price <= 2500
	case "over-2500":
		return product.Price > 2500
	default:
		return true
	}
}

func sortProducts(products []models.Product, sortMode string) {
	sort.SliceStable(products, func(left int, right int) bool {
		switch sortMode {
		case "latest":
			return products[left].UpdatedDaysAgo < products[right].UpdatedDaysAgo
		case "price-asc":
			return products[left].Price < products[right].Price
		case "price-desc":
			return products[left].Price > products[right].Price
		default:
			if products[left].Verified != products[right].Verified {
				return products[left].Verified
			}
			return products[left].Sales > products[right].Sales
		}
	})
}
