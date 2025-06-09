package implement

import (
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/request"
	"backend/internal/service"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type productServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) service.ProductService {
	return &productServiceImpl{
		productRepository: productRepository,
	}
}

func (s *productServiceImpl) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.productRepository.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả sản phẩm thất bại: %w", err)
	}

	return products, nil
}

func (s *productServiceImpl) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	product, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}

	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	return product, nil
}

func (s *productServiceImpl) CreateProduct(ctx context.Context, req request.CreateProductRequest) (*model.Product, error) {
	newProduct := &model.Product{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Brand:       req.Brand,
		Price:       req.Price,
		Inventory:   req.Inventory,
		Description: req.Description,
	}

	if err := s.productRepository.CreateProduct(ctx, newProduct); err != nil {
		return nil, fmt.Errorf("tạo sản phẩm thất bại: %w", err)
	}

	return newProduct, nil
}

func (s *productServiceImpl) UpdateProduct(ctx context.Context, id string, req *request.UpdateProductRequest) (*model.Product, error) {
	product, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}

	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	updateData := map[string]interface{}{}
	if req.Name != nil && *req.Name != product.Name {
		updateData["name"] = *req.Name
	}
	if req.Brand != nil && *req.Brand != product.Brand {
		updateData["brand"] = *req.Brand
	}
	if req.Price != nil && *req.Price != product.Price {
		updateData["price"] = *req.Price
	}
	if req.Inventory != nil && *req.Inventory != product.Inventory {
		updateData["inventory"] = *req.Inventory
	}
	if req.Description != nil && *req.Description != product.Description {
		updateData["description"] = *req.Description
	}

	if len(updateData) > 0 {
		if err = s.productRepository.UpdateProductByID(ctx, id, updateData); err != nil {
			if errors.Is(err, customErr.ErrProductNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật sản phẩm thất bại: %w", err)
		}
	}

	updatedProduct, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}

	if updatedProduct == nil {
		return nil, customErr.ErrProductNotFound
	}

	return updatedProduct, nil
}

func (s *productServiceImpl) DeleteProduct(ctx context.Context, id string) error {
	if err := s.productRepository.DeleteProductByID(ctx, id); err != nil {
		if errors.Is(err, customErr.ErrProductNotFound) {
			return err
		}
		return fmt.Errorf("xóa sản phẩm thất bại: %w", err)
	}

	return nil
}

func (s *productServiceImpl) DeleteManyProducts(ctx context.Context, req request.DeleteManyRequest) (int64, error) {
	productIDs := req.IDs

	rowsAccepted, err := s.productRepository.DeleteManyProducts(ctx, productIDs)
	if err != nil {
		return 0, fmt.Errorf("xóa danh sách sản phẩm thât bại: %w", err)
	}

	return rowsAccepted, nil
}
