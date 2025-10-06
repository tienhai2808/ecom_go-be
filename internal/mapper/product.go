package mapper

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/response"
)

func ToProductResponse(product *model.Product) *response.ProductResponse {
	return &response.ProductResponse{
		ID: product.ID,
		Name: product.Name,
		Slug: product.Slug,
		Description: product.Description,
		Price: product.Price,
		IsActive: product.IsActive,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
		Category: ToBaseCategoryResponse(product.Category),
		Inventory: ToInventoryResponse(product.Inventory),
		Images: ToImagesResponse(product.Images),
	}
}

func ToInventoryResponse(inv *model.Inventory) *response.InventoryResponse {
	return &response.InventoryResponse{
		ID: inv.ID,
		Quantity: inv.Quantity,
		Purchased: inv.Purchased,
		Stock: inv.Stock,
		IsStock: inv.IsStock,
	}
}

func ToImageResponse(img *model.Image) *response.ImageResponse {
	return &response.ImageResponse{
		ID:          img.ID,
		Url:         img.Url,
		IsThumbnail: img.IsThumbnail,
		SortOrder:   img.SortOrder,
	}
}

func ToImagesResponse(imgs []*model.Image) []*response.ImageResponse {
	var imgsResp []*response.ImageResponse
	for _, img := range imgs {
		imgsResp = append(imgsResp, ToImageResponse(img))
	}

	return imgsResp
}
