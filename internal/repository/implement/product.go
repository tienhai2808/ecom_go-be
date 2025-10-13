package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	internalType "github.com/tienhai2808/ecom_go/internal/types"
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

func (r *productRepositoryImpl) Search(ctx context.Context, query request.ProductPaginationQuery) (*internalType.ProductSearchResult, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	from := int((query.Page - 1) * query.Limit)
	size := int(query.Limit)

	esQuery := buildQuery(query)
	sort := buildSort(query)

	req := &search.Request{
		Query: esQuery,
		From:  &from,
		Size:  &size,
		Sort:  sort,
	}

	res, err := r.es.Search().
		Index("mysql_server.ecom_go.products").
		Request(req).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("tìm kiếm thất bại: %w", err)
	}

	productIDs := make([]int64, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		var doc struct {
			Payload struct {
				After struct {
					ID int64 `json:"id"`
				} `json:"after"`
			} `json:"payload"`
		}

		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			return nil, fmt.Errorf("lỗi giải mã document: %w", err)
		}
		productIDs = append(productIDs, doc.Payload.After.ID)
	}

	total := res.Hits.Total.Value
	totalPages := (total + int64(query.Limit) - 1) / int64(query.Limit)

	return &internalType.ProductSearchResult{
		IDs:        productIDs,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
		HasPrev:    query.Page > 1,
		HasNext:    int64(query.Page) < totalPages,
	}, nil
}

func (r *productRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Product, error) {
	return findByIDBase(ctx, r.db, id)
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

func (r *productRepositoryImpl) FindAllByIDWithThumbnail(ctx context.Context, ids []int64) ([]*model.Product, error) {
	var products []*model.Product
	if err := r.db.WithContext(ctx).Preload("Images").Scopes(getThumbnail).Where("id IN ?", ids).Find(&products).Error; err != nil {
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

func getThumbnail(db *gorm.DB) *gorm.DB {
	return db.Where("is_thumbnail = true")
}

func buildQuery(query request.ProductPaginationQuery) *types.Query {
	var mustQueries []types.Query

	if query.Search != "" {
		mustQueries = append(mustQueries, types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:  query.Search,
				Fields: []string{"payload.after.name^3", "payload.after.description"},
			},
		})
	}

	if query.IsActive != nil {
		activeValue := 0
		if *query.IsActive {
			activeValue = 1
		}
		mustQueries = append(mustQueries, types.Query{
			Term: map[string]types.TermQuery{
				"payload.after.is_active": {Value: activeValue},
			},
		})
	}

	if query.CategoryID != 0 {
		mustQueries = append(mustQueries, types.Query{
			Term: map[string]types.TermQuery{
				"payload.after.category_id": {Value: query.CategoryID},
			},
		})
	}

	if len(mustQueries) == 0 {
		return &types.Query{
			MatchAll: &types.MatchAllQuery{},
		}
	}

	if len(mustQueries) == 1 {
		return &mustQueries[0]
	}

	return &types.Query{
		Bool: &types.BoolQuery{
			Must: mustQueries,
		},
	}
}

func buildSort(query request.ProductPaginationQuery) []types.SortCombinations {
	if query.Sort == "" {
		if query.Search != "" {
			return []types.SortCombinations{
				types.SortOptions{
					SortOptions: map[string]types.FieldSort{
						"_score": {Order: &sortorder.Desc},
					},
				},
			}
		}

		return []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"payload.after.created_at": {Order: &sortorder.Desc},
				},
			},
		}
	}

	order := sortorder.Asc
	if query.Order == "desc" {
		order = sortorder.Desc
	}

	sortField := ""
	switch query.Sort {
	case "name":
		sortField = "payload.after.name.keyword"
	case "price":
		sortField = "payload.after.price"
	case "created_at":
		sortField = "payload.after.created_at"
	case "updated_at":
		sortField = "payload.after.updated_at"
	}

	return []types.SortCombinations{
		types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				sortField: {Order: &order},
			},
		},
	}
}
