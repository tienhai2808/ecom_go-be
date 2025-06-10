package handler

import (
	//customErr "backend/internal/errors"
	"backend/internal/imagekit"
	"backend/internal/service"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageHandler struct {
	imageService    service.ImageService
	imageKitService imagekit.ImageKitService
}

func NewImageHandler(imageService service.ImageService, imageKitService imagekit.ImageKitService) *ImageHandler {
	return &ImageHandler{
		imageService:    imageService,
		imageKitService: imageKitService,
	}
}

type ImageResponseDto struct {
	ImageURL   string
	ImageKitID string
}

func (h *ImageHandler) UploadImages(c *gin.Context) {
	ctx := c.Request.Context()
	form, err := c.MultipartForm()
	if err != nil {
		utils.JSON(c, http.StatusBadRequest, "Không thể parse muiltipart form", nil)
		return
	}
	defer form.RemoveAll()

	files := form.File["files"]
	if len(files) == 0 {
		utils.JSON(c, http.StatusBadRequest, "Không có hình ảnh nào được tải lên", nil)
		return
	}

	var imageResDtos []ImageResponseDto
	for _, file := range files {
		fileName := uuid.NewString()
		imageUrl, imageKitID, err := h.imageKitService.UploadImage(ctx, fileName, file)
		if err != nil {
			fmt.Printf("Lỗi up load ảnh %s ở ImageKitService %v\n", fileName, err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể upload ảnh lên ImageKit", nil)
			return
		}

		imageResDto := ImageResponseDto{
			ImageURL: imageUrl,
			ImageKitID: imageKitID,
		}

		imageResDtos = append(imageResDtos, imageResDto)
	}
	// if err != nil {
	// 	switch err {
	// 	case customErr.ErrProductNotFound:
	// 		utils.JSON(c, http.StatusNotFound, err.Error(), nil)
	// 	default:
	// 		fmt.Printf("Lỗi ở UploadImagesService: %v\n", err)
	// 		utils.JSON(c, http.StatusInternalServerError, "Không thể upload ảnh sản phẩm", nil)
	// 	}
	// 	return
	// }

	utils.JSON(c, http.StatusOK, "Upload ảnh sản phẩm thành công", gin.H{
		"images": imageResDtos,
	})
}
