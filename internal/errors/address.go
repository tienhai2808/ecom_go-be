package errors

import "errors"

var (
	ErrUserAddressNotFound = errors.New("không tìm thấy địa chỉ người dùng")

	ErrAddressNotFound = errors.New("địa chỉ không tồn tại")
	
	ErrExceedsQuantity = errors.New("vượt quá số lượng địa chỉ cho phép")
)