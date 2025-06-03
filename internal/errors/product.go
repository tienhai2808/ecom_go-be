package errors

import "errors"

var (
	ErrProductNotFound = errors.New("không tìm thấy sản phẩm")
	ErrProductExists   = errors.New("sản phẩm đã tồn tại")
)
