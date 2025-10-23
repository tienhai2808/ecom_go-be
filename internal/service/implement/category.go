package implement

import (
	"context"
	"errors"
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

func (s *categoryServiceImpl) UpdateCategory(ctx context.Context, id int64, req request.UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin danh mục sản phẩm thất bại: %w", err)
	}
	if category == nil {
		return nil, customErr.ErrCategoryNotFound
	}

	updateData := map[string]any{
		"name": req.Name,
		"slug": common.GenerateSlug(req.Name),
	}

	if err = s.categoryRepo.Update(ctx, id, updateData); err != nil {
		if common.IsUniqueViolation(err) {
			return nil, customErr.ErrCategorySlugAlreadyExists
		}
		if errors.Is(err, customErr.ErrCategoryNotFound) {
			return nil, customErr.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("cập nhật danh mục sản phẩm thất bại: %w", err)
	}

	category, err = s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin danh mục sản phẩm thất bại: %w", err)
	}
	if category == nil {
		return nil, customErr.ErrCategoryNotFound
	}

	return category, nil
}

func (s *categoryServiceImpl) DeleteCategory(ctx context.Context, id int64) error {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("lấy thông tin danh mục sản phẩm thất bại: %w", err)
	}
	if category == nil {
		return customErr.ErrCategoryNotFound
	}

	if err = s.categoryRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, customErr.ErrCategoryNotFound) {
			return err
		}
		return fmt.Errorf("xóa danh mục sản phẩm thất bại: %w", err)
	}

	return nil
}

func (s *categoryServiceImpl) DeleteCategories(ctx context.Context, req request.DeleteManyRequest) (int64, error) {
	categories, err := s.categoryRepo.FindAllByID(ctx, req.IDs)
	if err != nil {
		return 0, fmt.Errorf("lấy danh sách danh mục sản phẩm thât bại: %w", err)
	}
	if len(categories) != len(req.IDs) {
		return 0, customErr.ErrHasCategoryNotFound
	}

	rowsAccepted, err := s.categoryRepo.DeleteAllByID(ctx, req.IDs)
	if err != nil {
		return 0, fmt.Errorf("xóa danh sách danh mục sản phẩm thất bại: %w", err)
	}

	return rowsAccepted, nil
}
