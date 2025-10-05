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

func (r *imageRepositoryImpl) Create(ctx context.Context, image *model.Image) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *imageRepositoryImpl) Update(ctx context.Context, id int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Image{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrImageNotFound
	}

	return nil
}
