package implement

import (
	"context"
	"fmt"

	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
)

type categoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
	sfg          snowflake.SnowflakeGenerator
}

func NewCategoryService(categoryRepo repository.CategoryRepository, sfg snowflake.SnowflakeGenerator) service.CategoryService {
	return &categoryServiceImpl{
		categoryRepo,
		sfg,
	}
}

func (s *categoryServiceImpl) CreateCategory(ctx context.Context, req request.CreateCategoryRequest) (*model.Category, error) {
	slug := common.GenerateSlug(req.Name)

	categoryID, err := s.sfg.NextID()
	if err != nil {
		return nil, err
	}
	category := &model.Category{
		ID:   categoryID,
		Name: req.Name,
		Slug: slug,
	}
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		if common.IsUniqueViolation(err) {
			return nil, customErr.ErrCategorySlugAlreadyExists
		}
		return nil, fmt.Errorf("tạo danh mục sản phẩm thất bại: %w", err)
	}

	return category, nil
}

func (s *categoryServiceImpl) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách danh mục sản phẩm thất bại: %w", err)
	}

	return categories, nil
}