package implement

import (
	"context"
	"errors"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"

	"gorm.io/gorm"
)

type imageRepositoryImpl struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) repository.ImageRepository {
	return &imageRepositoryImpl{db}
}

func (r *imageRepositoryImpl) GetImageByID(ctx context.Context, id string) (*model.Image, error) {
	var image model.Image

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&image).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customErr.ErrImageNotFound
		}
		return nil, err
	}

	return &image, nil
}

func (r *imageRepositoryImpl) CreateImage(ctx context.Context, image *model.Image) error {
	if err := r.db.WithContext(ctx).Create(image).Error; err != nil {
		return err
	}

	return nil
}

func (r *imageRepositoryImpl) UpdateDownLoadUrlByID(ctx context.Context, id string, downloadUrl string) error {
	result := r.db.WithContext(ctx).Model(&model.Image{}).Where("id = ?", id).Update("download_url", downloadUrl)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrImageNotFound
	}

	return nil
}
