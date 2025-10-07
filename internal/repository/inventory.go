package repository

import (
	"context"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	UpdateTx(ctx context.Context, tx *gorm.DB, id int64, updateData map[string]any) error
}
