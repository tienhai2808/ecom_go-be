package implement

import (
	"backend/internal/model"
	"backend/internal/repository"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (r *productRepositoryImpl) GetAllProducts() ([]*model.Product, error) {
	var products []*model.Product

	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
