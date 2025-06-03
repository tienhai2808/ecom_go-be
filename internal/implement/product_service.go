package implement

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
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
	return s.productRepo.GetAllProducts()
}