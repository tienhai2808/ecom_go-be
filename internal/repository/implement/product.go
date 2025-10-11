package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
	es *elasticsearch.TypedClient
}

func NewProductRepository(db *gorm.DB, es *elasticsearch.TypedClient) repository.ProductRepository {
	return &productRepositoryImpl{
		db,
		es,
	}
}

func (r *productRepositoryImpl) FindAll(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	if err := r.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryImpl) Search(ctx context.Context, keyword string) ([]map[string]any, error) {
	res, err := r.es.Search().Index("mysql_server.ecom_go.products").Request(&search.Request{
		Query: &types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query: keyword,
				Fields: []string{"name^3", "description"},
	}}}).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("tìm kiếm thất bại: %w", err)
	}

	results := []map[string]any{}
	for _, hit := range res.Hits.Hits {
		var doc map[string]any
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			return nil, fmt.Errorf("lỗi giải mã document: %w", err)
		}
		results = append(results, doc)
	}

	return results, nil
}

func (r *productRepositoryImpl) FindByIDWithDetails(ctx context.Context, id int64) (*model.Product, error) {
	return findByIDBase(ctx, r.db, id, "Category", "Inventory", "Images")
}

func (r *productRepositoryImpl) FindByIDWithDetailsTx(ctx context.Context, tx *gorm.DB, id int64) (*model.Product, error) {
	return findByIDBase(ctx, tx, id, "Category", "Inventory", "Images")
}

func (r *productRepositoryImpl) FindByIDWithImages(ctx context.Context, id int64) (*model.Product, error) {
	return findByIDBase(ctx, r.db, id, "Images")
}

func (r *productRepositoryImpl) FindAllByIDWithImages(ctx context.Context, ids []int64) ([]*model.Product, error) {
	var products []*model.Product
	if err := r.db.WithContext(ctx).Preload("Images").Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepositoryImpl) UpdateTx(ctx context.Context, tx *gorm.DB, id int64, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *productRepositoryImpl) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Product{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrProductNotFound
	}

	return nil
}

func (r *productRepositoryImpl) DeleteAllByID(ctx context.Context, ids []int64) (int64, error) {
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Product{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func findByIDBase(ctx context.Context, tx *gorm.DB, id int64, preloads ...string) (*model.Product, error) {
	var product model.Product

	query := tx.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}
