package errors

import "errors"

var (
	ErrCategoryNotFound = errors.New("không tìm thấy danh mục sản phẩm")

	ErrCategorySlugAlreadyExists = errors.New("slug của danh mục sản phẩm đã tồn tại")

	ErrHasCategoryNotFound = errors.New("có danh mục sản phẩm không tìm thấy")
)