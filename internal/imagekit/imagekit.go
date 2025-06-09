package imagekit

import "mime/multipart"

type ImageKitService interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, string, error)
}