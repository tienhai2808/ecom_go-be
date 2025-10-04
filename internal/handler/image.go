package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/internal/service"
)

type ImageHandler struct {
	imageService    service.ImageService
}

func NewImageHandler(imageService service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService:    imageService,
	}
}

func (h *ImageHandler) UploadImages(c *gin.Context) {

}