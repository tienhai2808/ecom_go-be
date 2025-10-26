package errors

import "errors"

var (
	ErrKeyNotFound = errors.New("key đã hết hạn hoặc không tìm thấy")

	ErrUnauthorized = errors.New("truy cập không hợp lệ")

	ErrInvalidID = errors.New("ID không hợp lệ")

	ErrInvalidRequest = errors.New("yêu cầu không hợp lệ")

	ErrGuestIdNotFound = errors.New("không tìm thấy guest_id")
)