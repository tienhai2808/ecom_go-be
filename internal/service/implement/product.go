package implement

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"fmt"
)

type productServiceImpl struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) service.ProductService {
	return &productServiceImpl{
		productRepo: productRepo,
	}
}

func (s *productServiceImpl) GetAllProducts() ([]*model.Product, error) {
	products, err := s.productRepo.GetAllProducts()
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả sản phẩm thất bại: %w", err)
	}

	return products, nil
}