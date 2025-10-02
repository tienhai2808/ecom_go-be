package implement

import (
	"context"
	"errors"
	"fmt"

	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/util"
)

type productServiceImpl struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
}

func NewProductService(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository) service.ProductService {
	return &productServiceImpl{
		productRepository,
		categoryRepository,
	}
}

func (s *productServiceImpl) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.productRepository.GetAllProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả sản phẩm thất bại: %w", err)
	}

	return products, nil
}

func (s *productServiceImpl) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
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
	productID, err := util.NewSnowflakeID()
	if err != nil {
		return nil, err
	}
	inventoryID, err := util.NewSnowflakeID()
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepository.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin danh mục sản phẩm thất bại: %w", err)
	}
	if category == nil {
		return nil, customErr.ErrCategoryNotFound
	}

	newProduct := &model.Product{
		ID:          productID,
		Name:        req.Name,
		Price:       req.Price,
		Slug:        util.GenerateSlug(req.Name),
		CategoryID:  category.ID,
		Description: req.Description,
		Inventory: &model.Inventory{
			ID:        inventoryID,
			Quantity:  req.Quantity,
			Purchased: 0,
		},
	}
	newProduct.Inventory.SetStock()

	if err := s.productRepository.CreateProduct(ctx, newProduct); err != nil {
		return nil, fmt.Errorf("tạo sản phẩm thất bại: %w", err)
	}

	return newProduct, nil
}

func (s *productServiceImpl) UpdateProduct(ctx context.Context, id int64, req *request.UpdateProductRequest) (*model.Product, error) {
	product, err := s.productRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}
	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	updateData := map[string]any{}
	if req.Name != nil && *req.Name != product.Name {
		updateData["name"] = *req.Name
	}
	if req.CategoryID != nil && *req.CategoryID != product.CategoryID {
		updateData["category_id"] = *req.CategoryID
	}
	if req.Price != nil && *req.Price != product.Price {
		updateData["price"] = *req.Price
	}
	if req.Quantity != nil && *req.Quantity != product.Inventory.Quantity {

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

func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
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
