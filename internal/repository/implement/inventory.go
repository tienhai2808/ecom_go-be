package implement

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"gorm.io/gorm"
)

type inventoryRepositoryImpl struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) repository.InventoryRepository {
	return &inventoryRepositoryImpl{db}
}

func (r *inventoryRepositoryImpl) UpdateTx(ctx context.Context, tx *gorm.DB, id int64, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.Inventory{}).Where("id = ?", id).Updates(updateData).Error
}
