package implement

import (
	"context"
	"errors"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepositoryImpl{db}
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *categoryRepositoryImpl) FindAll(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	return r.FindByIDTx(ctx, r.db, id)
}

func (r *categoryRepositoryImpl) FindByIDTx(ctx context.Context, tx *gorm.DB, id int64) (*model.Category, error) {
	var category model.Category
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}
