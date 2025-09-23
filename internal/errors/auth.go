package errors

import "errors"

var (
	ErrUnAuth = errors.New("người dùng chưa đăng nhập")

	ErrInvalidToken = errors.New("token không hợp lệ hoặc đã hết hạn")
)