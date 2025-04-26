package auth

import "fmt"

var (
	ErrUsernameExists    = fmt.Errorf("username đã tồn tại")
	ErrEmailExists       = fmt.Errorf("email đã tồn tại")
	ErrDatabase          = fmt.Errorf("không thể kết nối đến cơ sở dữ liệu")
	ErrInvalidOTP        = fmt.Errorf("OTP không hợp lệ")
	ErrTokenExpired      = fmt.Errorf("token đã hết hạn hoặc không tồn tại")
	ErrTooManyAttempts   = fmt.Errorf("vượt quá số lần thử OTP")
	ErrUsernameNotFound  = fmt.Errorf("người dùng không tồn tại")
	ErrIncorrectPassword = fmt.Errorf("mật khẩu không chính xác")
)
