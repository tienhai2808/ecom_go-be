package errors

import "errors"

var (
	ErrKeyNotFound = errors.New("key đã hết hạn hoặc không tìm thấy")
	ErrUpdateFailed = errors.New("cập nhật thất bại")
)