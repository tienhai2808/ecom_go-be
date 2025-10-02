package handler

import (
	//customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"fmt"
	"github.com/tienhai2808/ecom_go/internal/imagekit"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/util"
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
		util.JSON(c, http.StatusBadRequest, "Không thể parse muiltipart form", nil)
		return
	}
	defer form.RemoveAll()

	files := form.File["files"]
	if len(files) == 0 {
		util.JSON(c, http.StatusBadRequest, "Không có hình ảnh nào được tải lên", nil)
		return
	}

	var imageResDtos []ImageResponseDto
	for _, file := range files {
		fileName := uuid.NewString()
		imageUrl, imageKitID, err := h.imageKitService.UploadImage(ctx, fileName, file)
		if err != nil {
			fmt.Printf("Lỗi up load ảnh %s ở ImageKitService %v\n", fileName, err)
			util.JSON(c, http.StatusInternalServerError, "Không thể upload ảnh lên ImageKit", nil)
			return
		}

		imageResDto := ImageResponseDto{
			ImageURL:   imageUrl,
			ImageKitID: imageKitID,
		}

		imageResDtos = append(imageResDtos, imageResDto)
	}
	// if err != nil {
	// 	switch err {
	// 	case customErr.ErrProductNotFound:
	// 		util.JSON(c, http.StatusNotFound, err.Error(), nil)
	// 	default:
	// 		fmt.Printf("Lỗi ở UploadImagesService: %v\n", err)
	// 		util.JSON(c, http.StatusInternalServerError, "Không thể upload ảnh sản phẩm", nil)
	// 	}
	// 	return
	// }

	util.JSON(c, http.StatusOK, "Upload ảnh sản phẩm thành công", gin.H{
		"images": imageResDtos,
	})
}
