package auth

import "errors"

var (
	ErrUsernameExists    = errors.New("username đã tồn tại")
	ErrEmailExists       = errors.New("email đã tồn tại")
	ErrDatabase          = errors.New("không thể kết nối đến cơ sở dữ liệu")
	ErrInvalidOTP        = errors.New("OTP không hợp lệ")
	ErrTokenExpired      = errors.New("token đã hết hạn hoặc không tồn tại")
	ErrTooManyAttempts   = errors.New("vượt quá số lần thử OTP")
	ErrUserNotFound      = errors.New("người dùng không tồn tại")
	ErrIncorrectPassword = errors.New("mật khẩu không chính xác")
	ErrUnAuth            = errors.New("không có quyền truy cập")
	ErrUpdateFailed      = errors.New("cập nhật thông tin người dùng thất bại")
)
