package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/types"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]*model.Product, error)

	Search(ctx context.Context, query request.ProductPaginationQuery) (*types.ProductSearchResult, error)

	FindByIDWithDetails(ctx context.Context, id int64) (*model.Product, error)

	FindByID(ctx context.Context, id int64) (*model.Product, error)

	FindByIDWithImages(ctx context.Context, id int64) (*model.Product, error)

	FindByIDWithDetailsTx(ctx context.Context, tx *gorm.DB, id int64) (*model.Product, error)

	FindAllByIDWithImages(ctx context.Context, ids []int64) ([]*model.Product, error)

	FindAllByIDWithThumbnail(ctx context.Context, ids []int64) ([]*model.Product, error)

	FindAllByIDWithCategoryAndThumbnail(ctx context.Context, ids []int64) ([]*model.Product, error)

	Create(ctx context.Context, product *model.Product) error

	UpdateTx(ctx context.Context, tx *gorm.DB, id int64, updateData map[string]any) error

	Delete(ctx context.Context, id int64) error

	DeleteAllByID(ctx context.Context, ids []int64) (int64, error)
}
