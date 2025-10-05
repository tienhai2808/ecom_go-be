package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"net/http"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/rabbitmq"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"github.com/tienhai2808/ecom_go/internal/types"
)

type productServiceImpl struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	rabbitChan  *amqp091.Channel
	sfg          snowflake.SnowflakeGenerator
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, rabbitChan  *amqp091.Channel, sfg snowflake.SnowflakeGenerator) service.ProductService {
	return &productServiceImpl{
		productRepo,
		categoryRepo,
		rabbitChan,
		sfg,
	}
}

func (s *productServiceImpl) GetAllProducts(ctx context.Context) ([]*model.Product, error) {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả sản phẩm thất bại: %w", err)
	}

	return products, nil
}

func (s *productServiceImpl) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}

	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	return product, nil
}

func (s *productServiceImpl) CreateProduct(ctx context.Context, req *request.CreateProductForm) (*model.Product, error) {
	productID, err := s.sfg.NextID()
	if err != nil {
		return nil, err
	}
	inventoryID, err := s.sfg.NextID()
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepo.FindByID(ctx, req.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin danh mục sản phẩm thất bại: %w", err)
	}
	if category == nil {
		return nil, customErr.ErrCategoryNotFound
	}

	slug := common.GenerateSlug(req.Name)

	imgQuan := len(req.Images)
	images := make([]*model.Image, imgQuan)
	publishCh := make(chan *types.UploadImageMessage, imgQuan)

	if imgQuan > 0 {
		for _, img := range req.Images {
			fileName := fmt.Sprintf("%s_%d%s", slug, img.SortOrder, getFileExt(img.File))
			imageID, err := s.sfg.NextID()
			if err != nil {
				return nil, err
			}

			newImg := &model.Image{
				ID: imageID,
				IsThumbnail: *img.IsThumbnail,
				SortOrder: img.SortOrder,
			}

			uploadReq := &types.UploadImageMessage{
				ImageID: imageID,
				FileName: fileName,
				FileData: img.File,
			}

			publishCh <- uploadReq
			images = append(images, newImg)
		}
	}
	close(publishCh)

	newProduct := &model.Product{
		ID:          productID,
		Name:        req.Name,
		Price:       req.Price,
		Slug:        slug,
		CategoryID:  category.ID,
		Description: req.Description,
		Inventory: &model.Inventory{
			ID:        inventoryID,
			Quantity:  req.Quantity,
			Purchased: 0,
		},
		Images: images,
	}
	newProduct.Inventory.SetStock()

	if err := s.productRepo.Create(ctx, newProduct); err != nil {
		return nil, fmt.Errorf("tạo sản phẩm thất bại: %w", err)
	}

	go func ()  {
		for req := range publishCh {
			body, _ := json.Marshal(req)
			if err := rabbitmq.PublishMessage(s.rabbitChan, common.ExchangeImage, common.RoutingKeyImageUpload, body); err != nil {
				log.Printf("đẩy tin nhắn upload ảnh thất bại: %v", err)
			}
		}
	}()

	return newProduct, nil
}

func (s *productServiceImpl) UpdateProduct(ctx context.Context, id int64, req *request.UpdateProductRequest) (*model.Product, error) {
	product, err := s.productRepo.FindByID(ctx, id)
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
		if err = s.productRepo.Update(ctx, id, updateData); err != nil {
			if errors.Is(err, customErr.ErrProductNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật sản phẩm thất bại: %w", err)
		}
	}

	updatedProduct, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}

	if updatedProduct == nil {
		return nil, customErr.ErrProductNotFound
	}

	return updatedProduct, nil
}

func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.productRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, customErr.ErrProductNotFound) {
			return err
		}
		return fmt.Errorf("xóa sản phẩm thất bại: %w", err)
	}

	return nil
}

func (s *productServiceImpl) DeleteManyProducts(ctx context.Context, req request.DeleteManyRequest) (int64, error) {
	productIDs := req.IDs
	rowsAccepted, err := s.productRepo.DeleteAllByID(ctx, productIDs)
	if err != nil {
		return 0, fmt.Errorf("xóa danh sách sản phẩm thât bại: %w", err)
	}

	return rowsAccepted, nil
}

func getFileExt(data []byte) string {
	contentType := http.DetectContentType(data)

	exts, _ := mime.ExtensionsByType(contentType)
	if len(exts) > 0 {
		return exts[0]
	}

	return ".webp"
}
