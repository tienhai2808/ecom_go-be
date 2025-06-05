package errors

import "errors"

var (
	ErrKeyNotFound = errors.New("key đã hết hạn hoặc không tìm thấy")
	ErrUnauthorized = errors.New("truy cập không hợp lệ")
)