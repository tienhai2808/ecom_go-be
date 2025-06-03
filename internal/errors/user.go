package errors

import "errors"

var (
	ErrEmailExists       = errors.New("email đã tồn tại")
	ErrUsernameExists    = errors.New("username đã tồn tại")
	ErrTooManyAttempts   = errors.New("vượt quá số lần thử OTP")
	ErrInvalidOTP        = errors.New("OTP không hợp lệ")
	ErrUserNotFound      = errors.New("không tìm thấy người dùng")
	ErrIncorrectPassword = errors.New("mật khẩu không chính xác")
)
