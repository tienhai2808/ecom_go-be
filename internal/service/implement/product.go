package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

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
	"gorm.io/gorm"
)

type productServiceImpl struct {
	productRepo   repository.ProductRepository
	categoryRepo  repository.CategoryRepository
	inventoryRepo repository.InventoryRepository
	imageRepo     repository.ImageRepository
	db            *gorm.DB
	rabbitChan    *amqp091.Channel
	sfg           snowflake.SnowflakeGenerator
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, inventoryRepo repository.InventoryRepository, imageRepo repository.ImageRepository, db *gorm.DB, rabbitChan *amqp091.Channel, sfg snowflake.SnowflakeGenerator) service.ProductService {
	return &productServiceImpl{
		productRepo,
		categoryRepo,
		inventoryRepo,
		imageRepo,
		db,
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
	product, err := s.productRepo.FindByIDWithDetails(ctx, id)
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
	images := make([]*model.Image, 0, imgQuan)
	publishCh := make(chan *types.UploadImageMessage, imgQuan)

	if imgQuan > 0 {
		for _, img := range req.Images {
			fileName := fmt.Sprintf("%s_%d", slug, img.SortOrder)
			imageID, err := s.sfg.NextID()
			if err != nil {
				return nil, err
			}

			newImg := &model.Image{
				ID:          imageID,
				IsThumbnail: *img.IsThumbnail,
				SortOrder:   img.SortOrder,
			}

			uploadReq := &types.UploadImageMessage{
				ImageID:  imageID,
				FileName: fileName,
				FileData: img.FileData,
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
		Category:    category,
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
		if common.IsUniqueViolation(err) {
			return nil, customErr.ErrProductSlugAlreadyExists
		}
		return nil, fmt.Errorf("tạo sản phẩm thất bại: %w", err)
	}

	go func() {
		for req := range publishCh {
			body, _ := json.Marshal(req)
			if err := rabbitmq.PublishMessage(s.rabbitChan, common.ExchangeImage, common.RoutingKeyImageUpload, body); err != nil {
				log.Printf("đẩy tin nhắn upload ảnh thất bại: %v", err)
			}
		}
	}()

	return newProduct, nil
}

func (s *productServiceImpl) UpdateProduct(ctx context.Context, id int64, req *request.UpdateProductForm) (*model.Product, error) {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		product, err := s.productRepo.FindByIDWithDetailsTx(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
		}
		if product == nil {
			return customErr.ErrProductNotFound
		}

		updateData := map[string]any{}
		if req.Name != nil && *req.Name != product.Name {
			updateData["name"] = *req.Name
			updateData["slug"] = common.GenerateSlug(*req.Name)
		}
		if req.Price != nil && *req.Price != product.Price {
			updateData["price"] = *req.Price
		}
		if req.Description != nil && *req.Description != product.Description {
			updateData["description"] = *req.Description
		}
		if req.IsActive != nil && *req.IsActive != product.IsActive {
			updateData["is_active"] = *req.IsActive
		}

		if req.CategoryID != nil && *req.CategoryID != product.CategoryID {
			category, err := s.categoryRepo.FindByIDTx(ctx, tx, *req.CategoryID)
			if err != nil {
				return fmt.Errorf("lấy thông tin danh mục sản phẩm thất bại: %w", err)
			}
			if category == nil {
				return customErr.ErrCategoryNotFound
			}

			updateData["category_id"] = category.ID
		}

		if len(updateData) > 0 {
			if err = s.productRepo.UpdateTx(ctx, tx, id, updateData); err != nil {
				if common.IsUniqueViolation(err) {
					return customErr.ErrProductSlugAlreadyExists
				}
				return fmt.Errorf("cập nhật thông tin sản phẩm thất bại: %w", err)
			}
		}

		if req.Quantity != nil && *req.Quantity != product.Inventory.Quantity && *req.Quantity >= product.Inventory.Purchased {
			updateData := map[string]any{
				"quantity": *req.Quantity,
				"stock":    gorm.Expr("quantity - purchased"),
				"is_stock": gorm.Expr("CASE WHEN (quantity - purchased) <= 5 THEN false ELSE true END"),
			}

			if err = s.inventoryRepo.UpdateTx(ctx, tx, product.Inventory.ID, updateData); err != nil {
				return fmt.Errorf("cập nhật số lượng sản phẩm thất bại: %w", err)
			}
		}

		if len(req.DeleteImageIDs) > 0 {
			imgs, err := s.imageRepo.FindAllByIDTx(ctx, tx, req.DeleteImageIDs)
			if err != nil {
				return fmt.Errorf("lấy danh sách hình ảnh sản phẩm xóa thất bại: %w", err)
			}
			if len(imgs) != len(req.DeleteImageIDs) {
				return customErr.ErrHasImageNotFound
			}

			if err := s.imageRepo.DeleteAllByIDTx(ctx, tx, req.DeleteImageIDs); err != nil {
				return fmt.Errorf("xóa danh sách hình ảnh sản phẩm thất bại: %w", err)
			}

			publishChan := make(chan string, len(imgs))
			for _, img := range imgs {
				if strings.TrimSpace(img.PublicID) != "" {
					publishChan <- img.PublicID
				}
			}
			close(publishChan)

			go func() {
				for req := range publishChan {
					body := []byte(req)
					if err := rabbitmq.PublishMessage(s.rabbitChan, common.ExchangeImage, common.RoutingKeyImageDelete, body); err != nil {
						log.Printf("đẩy tin nhắn xóa ảnh thất bại: %v", err)
					}
				}
			}()
		}

		if len(req.UpdateImages) > 0 {
			imgIDs := make([]int64, 0, len(req.UpdateImages))
			for _, img := range req.UpdateImages {
				imgIDs = append(imgIDs, img.ID)
			}

			imgs, err := s.imageRepo.FindAllByIDTx(ctx, tx, imgIDs)
			if err != nil {
				return fmt.Errorf("lấy danh sách chỉnh sửa hình ảnh sản phẩm thất bại: %w", err)
			}
			if len(imgIDs) != len(imgs) {
				return customErr.ErrHasImageNotFound
			}

			for _, img := range req.UpdateImages {
				updateData := map[string]any{}
				if img.IsThumbnail != nil {
					updateData["is_thumbnail"] = img.IsThumbnail
				}
				if img.SortOrder != nil {
					updateData["sort_order"] = img.SortOrder
				}

				if len(updateData) > 0 {
					if err = s.imageRepo.UpdateTx(ctx, tx, img.ID, updateData); err != nil {
						return fmt.Errorf("cập nhật hình ảnh thất bại: %w", err)
					}
				}
			}
		}

		imgQuan := len(req.NewImages)
		if imgQuan > 0 {
			images := make([]*model.Image, 0, imgQuan)
			publishCh := make(chan *types.UploadImageMessage, imgQuan)

			for _, img := range req.NewImages {
				var slug string
				if req.Name != nil {
					slug = common.GenerateSlug(*req.Name)
				} else {
					slug = product.Slug
				}

				fileName := fmt.Sprintf("%s_%d", slug, img.SortOrder)
				imageID, err := s.sfg.NextID()
				if err != nil {
					return err
				}

				fmt.Println(*img.IsThumbnail)

				newImg := &model.Image{
					ID:          imageID,
					IsThumbnail: *img.IsThumbnail,
					SortOrder:   img.SortOrder,
					ProductID:   product.ID,
				}

				uploadReq := &types.UploadImageMessage{
					ImageID:  imageID,
					FileName: fileName,
					FileData: img.FileData,
				}

				publishCh <- uploadReq
				images = append(images, newImg)
			}
			close(publishCh)

			if err = s.imageRepo.CreateAllTx(ctx, tx, images); err != nil {
				return fmt.Errorf("tạo hình ảnh thất bại: %w", err)
			}

			go func() {
				for req := range publishCh {
					body, _ := json.Marshal(req)
					if err := rabbitmq.PublishMessage(s.rabbitChan, common.ExchangeImage, common.RoutingKeyImageUpload, body); err != nil {
						log.Printf("đẩy tin nhắn upload ảnh thất bại: %v", err)
					}
				}
			}()
		}

		return nil
	}); err != nil {
		return nil, err
	}

	updatedProduct, err := s.productRepo.FindByIDWithDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}
	if updatedProduct == nil {
		return nil, customErr.ErrProductNotFound
	}

	return updatedProduct, nil
}

func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	product, err := s.productRepo.FindByIDWithImages(ctx, id)
	if err != nil {
		return fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}
	if product == nil {
		return customErr.ErrProductNotFound
	}

	if err := s.productRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, customErr.ErrProductNotFound) {
			return err
		}
		return fmt.Errorf("xóa sản phẩm thất bại: %w", err)
	}

	if len(product.Images) > 0 {
		publishChan := make(chan string, len(product.Images))
		for _, img := range product.Images {
			if strings.TrimSpace(img.PublicID) != "" {
				publishChan <- img.PublicID
			}
		}
		close(publishChan)

		go func() {
			for req := range publishChan {
				body := []byte(req)
				if err := rabbitmq.PublishMessage(s.rabbitChan, common.ExchangeImage, common.RoutingKeyImageDelete, body); err != nil {
					log.Printf("đẩy tin nhắn xóa ảnh thất bại: %v", err)
				}
			}
		}()
	}

	return nil
}

func (s *productServiceImpl) DeleteProducts(ctx context.Context, req request.DeleteManyRequest) (int64, error) {
	products, err := s.productRepo.FindAllByIDWithImages(ctx, req.IDs)
	if err != nil {
		return 0, fmt.Errorf("lấy danh sách sản phẩm cần xóa thất bại: %w", err)
	}
	if len(req.IDs) != len(products) {
		return 0, customErr.ErrHasProductNotFound
	}

	rowsAccepted, err := s.productRepo.DeleteAllByID(ctx, req.IDs)
	if err != nil {
		return 0, fmt.Errorf("xóa danh sách sản phẩm thât bại: %w", err)
	}

	imgPublicIDs := []string{}
	seen := make(map[string]bool)
	for _, product := range products {
		for _, img := range product.Images {
			if strings.TrimSpace(img.PublicID) != "" && !seen[img.PublicID] {
				seen[img.PublicID] = true
				imgPublicIDs = append(imgPublicIDs, img.PublicID)
			}
		}
	}

	publishChan := make(chan string, len(imgPublicIDs))
	for _, publicID := range imgPublicIDs {
		publishChan <- publicID
	}
	close(publishChan)

	go func ()  {
		for req := range publishChan {
			body := []byte(req)
			if err := rabbitmq.PublishMessage(s.rabbitChan, common.ExchangeImage, common.RoutingKeyImageDelete, body); err != nil {
				log.Printf("đẩy tin nhắn xóa ảnh thất bại: %v", err)
			}
		}
	}()

	return rowsAccepted, nil
}
