package consumers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/tienhai2808/ecom_go/internal/cloudinary"
	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/initialization"
	"github.com/tienhai2808/ecom_go/internal/rabbitmq"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/types"
)

func StartDeleteImageMessage(mqc *initialization.RabbitMQConn, cld cloudinary.CloudinaryService) {
	if err := rabbitmq.ConsumeMessage(mqc.Chan, common.QueueNameImageDelete, common.ExchangeImage, common.RoutingKeyImageDelete, func(body []byte) error {
		publicID := string(body)
		ctx := context.Background()

		if err := cld.DeleteFile(ctx, publicID, "image"); err != nil {
			return fmt.Errorf("xóa file thất bại: %w", err)
		}
		log.Printf("Xóa hình ảnh có PublicID: %s thành công", publicID)

		return nil
	}); err != nil {
		log.Printf("Lỗi khởi tạo delete image consumer: %v", err)
	}
}

func StartUploadImageMessage(mqc *initialization.RabbitMQConn, cld cloudinary.CloudinaryService, imageRepo repository.ImageRepository) {
	if err := rabbitmq.ConsumeMessage(mqc.Chan, common.QueueNameImageUpload, common.ExchangeImage, common.RoutingKeyImageUpload, func(body []byte) error {
		var msg types.UploadImageMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return fmt.Errorf("chuyển đổi tin nhắn upload ảnh thất bại: %w", err)
		}

		ctx := context.Background()

		res, err := cld.UploadBinaryFile(ctx, msg.FileData, msg.FileName)
		if err != nil {
			return fmt.Errorf("upload ảnh thất bại: %w", err)
		}
		log.Printf("Upload ảnh %s lên Cloudinary thành công", res.URL)

		updateData := map[string]any{
			"public_id": res.PublicID,
			"url":       res.URL,
		}

		if err = imageRepo.Update(ctx, msg.ImageID, updateData); err != nil {
			if errors.Is(err, customErr.ErrImageNotFound) {
				return err
			}
			return fmt.Errorf("cập nhật hình ảnh thất bại: %w", err)
		}
		log.Printf("Cập nhật ảnh có ID %d thành công", msg.ImageID)

		return nil
	}); err != nil {
		log.Printf("Lỗi khởi tạo upload image consumer: %v", err)
	}
}
