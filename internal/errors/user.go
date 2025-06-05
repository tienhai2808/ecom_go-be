package errors

import "errors"

var (
	ErrEmailExists       = errors.New("email đã tồn tại")
	ErrUsernameExists    = errors.New("username đã tồn tại")
	ErrTooManyAttempts   = errors.New("vượt quá số lần thử OTP")
	ErrInvalidOTP        = errors.New("OTP không hợp lệ")
	ErrUserNotFound      = errors.New("người dùng không tồn tại")
	ErrIncorrectPassword = errors.New("mật khẩu không chính xác")
	ErrUserConflict      = errors.New("không thể xóa chính tài khoản của bạn")
)
