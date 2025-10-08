package errors

import "errors"

var (
	ErrProductNotFound = errors.New("không tìm thấy sản phẩm")

	ErrProductExists   = errors.New("sản phẩm đã tồn tại")

	ErrProductSlugAlreadyExists = errors.New("slug của sản phẩm đã tồn tại")

	ErrHasProductNotFound = errors.New("có sản phẩm không tìm thấy")
)
