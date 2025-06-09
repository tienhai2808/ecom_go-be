package implement

import (
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"context"
	"fmt"
	"mime/multipart"
)

type imageServiceImpl struct {
	imageRepository   repository.ImageRepository
	productRepository repository.ProductRepository
}

func NewImageService(imageRepository repository.ImageRepository, productRepository repository.ProductRepository) service.ImageService {
	return &imageServiceImpl{
		imageRepository:   imageRepository,
		productRepository: productRepository,
	}
}

func (s *imageServiceImpl) UploadImages(ctx context.Context, files []*multipart.FileHeader, productID string) ([]*model.Image, error) {
	product, err := s.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}

	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	return nil, nil
}
