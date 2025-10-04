package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/internal/service"
)

type ImageHandler struct {
	imageSvc service.ImageService
}

func NewImageHandler(imageSvc service.ImageService) *ImageHandler {
	return &ImageHandler{imageSvc}
}

func (h *ImageHandler) UploadImages(c *gin.Context) {}
