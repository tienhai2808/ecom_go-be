package handler

import (
	customErr "backend/internal/errors"
	"backend/internal/imagekit"
	"backend/internal/service"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageHandler struct {
	imageService service.ImageService
	imageKitService imagekit.ImageKitService
}

func NewImageHandler(imageService service.ImageService, imageKitService imagekit.ImageKitService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
		imageKitService: imageKitService,
	}
}

func (h *ImageHandler) UploadImages(c *gin.Context) {
	ctx := c.Request.Context()
	form, err := c.MultipartForm()
	if err != nil {
		utils.JSON(c, http.StatusBadRequest, "Không thể parse muiltipart form", nil)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.JSON(c, http.StatusBadRequest, "Không có hình ảnh nào được tải lên", nil)
		return
	}

	productID := c.PostForm("product_id")

	uploadedImages, err := h.imageService.UploadImages(ctx, files, productID)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			utils.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở UploadImagesService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể upload ảnh sản phẩm", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Upload ảnh sản phẩm thành công", gin.H{
		"images": uploadedImages,
	})
}
