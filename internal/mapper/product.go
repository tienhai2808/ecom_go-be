package mapper

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/response"
)

func ToProductResponse(product *model.Product) *response.ProductResponse {
	return &response.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Category:    ToBaseCategoryResponse(product.Category),
		Inventory:   ToInventoryResponse(product.Inventory),
		Images:      ToImagesResponse(product.Images),
	}
}

func ToBaseProductResponse(product *model.Product) *response.BaseProductResponse {
	return &response.BaseProductResponse{
		ID: product.ID,
		Name: product.Name,
		Slug: product.Slug,
		Price: product.Price,
		IsActive: product.IsActive,
		Thumbnail: product.Images[0].Url,
	}
}

func ToBaseProductsResponse(prds []*model.Product) []*response.BaseProductResponse {
	if len(prds) == 0 {
		return make([]*response.BaseProductResponse, 0)
	}

	prdsResp := make([]*response.BaseProductResponse, 0, len(prds)) 
	for _, prd := range prds {
		prdsResp = append(prdsResp, ToBaseProductResponse(prd))
	}

	return prdsResp
}

func ToInventoryResponse(inv *model.Inventory) *response.InventoryResponse {
	return &response.InventoryResponse{
		ID:        inv.ID,
		Quantity:  inv.Quantity,
		Purchased: inv.Purchased,
		Stock:     inv.Stock,
		IsStock:   inv.IsStock,
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
	if len(imgs) == 0 {
		return make([]*response.ImageResponse, 0)
	}

	imgsResp := make([]*response.ImageResponse, 0, len(imgs))
	for _, img := range imgs {
		imgsResp = append(imgsResp, ToImageResponse(img))
	}

	return imgsResp
}
