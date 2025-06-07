package implement

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/request"
	"backend/internal/service"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type productServiceImpl struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) service.ProductService {
	return &productServiceImpl{
		productRepo: productRepo,
	}
}

func (s *productServiceImpl) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.productRepo.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả sản phẩm thất bại: %w", err)
	}

	return products, nil
}

func (s *productServiceImpl) CreateProduct(ctx context.Context, req request.CreateProductRequest) (*model.Product, error) {
	newProduct := &model.Product{
		ID: uuid.NewString(),
		Name: req.Name,
		Brand: req.Brand,
		Price: req.Price,
		Inventory: req.Inventory,
		Description: req.Description,
	}

	if err := s.productRepo.CreateProduct(ctx, newProduct); err != nil {
		return nil, fmt.Errorf("tạo sản phẩm thất bại: %w", err)
	}

	return newProduct, nil
}