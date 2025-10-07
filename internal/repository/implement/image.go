package implement

import (
	"context"
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

func (r *imageRepositoryImpl) Create(ctx context.Context, image *model.Image) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *imageRepositoryImpl) CreateAllTx(ctx context.Context, tx *gorm.DB, images []*model.Image) error {
	return tx.WithContext(ctx).Create(images).Error
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

func (r *imageRepositoryImpl) UpdateTx(ctx context.Context, tx *gorm.DB, id int64, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.Image{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *imageRepositoryImpl) FindAllByIDTx(ctx context.Context, tx *gorm.DB, ids []int64) ([]*model.Image, error) {
	var images []*model.Image
	if err := tx.WithContext(ctx).Where("id IN ?", ids).Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (r *imageRepositoryImpl) DeleteAllByIDTx(ctx context.Context, tx *gorm.DB, ids []int64) error {
	return tx.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Image{}).Error
}