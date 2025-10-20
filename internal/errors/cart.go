package errors

import "errors"

var (
	ErrCartNotFound = errors.New("không tìm thấy giỏ hàng")

	ErrCartItemNotFound = errors.New("không tìm thấy item trong giỏ hàng")
)