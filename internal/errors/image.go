package errors

import "errors"

var (
	ErrImageNotFound = errors.New("không tìm thấy hình ảnh")

	ErrHasImageNotFound = errors.New("có hình ảnh không tìm thấy")
)