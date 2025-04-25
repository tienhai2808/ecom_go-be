package auth

import (
	"errors"
	"fmt"
)

var (
	ErrUsernameExists = errors.New("username đã tồn tại")
	ErrEmailExists    = errors.New("email đã tồn tại")
	ErrDatabase       = errors.New("không thể kết nối đến cơ sở dữ liệu")
	ErrInvalidOTP     = fmt.Errorf("OTP không hợp lệ")
	ErrTokenExpired   = fmt.Errorf("token đã hết hạn hoặc không tồn tại")
	ErrTooManyAttempts = fmt.Errorf("vượt quá số lần thử OTP")
) 