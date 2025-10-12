package response

import "time"

type ProductResponse struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Slug        string                `json:"slug"`
	Price       float64               `json:"price"`
	Description string                `json:"description"`
	IsActive    bool                  `json:"is_active"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Category    *BaseCategoryResponse `json:"category"`
	Inventory   *InventoryResponse    `json:"inventory"`
	Images      []*ImageResponse      `json:"images"`
}

type ImageResponse struct {
	ID          int64  `json:"id"`
	Url         string `json:"url"`
	IsThumbnail bool   `json:"is_thumbnail"`
	SortOrder   int    `json:"sort_order"`
}

type InventoryResponse struct {
	ID        int64 `json:"id"`
	Quantity  uint  `json:"quantity"`
	Purchased uint  `json:"purchased"`
	Stock     uint  `json:"stock"`
	IsStock   bool  `json:"is_stock"`
}

type BaseProductResponse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Price     float64 `json:"price"`
	IsActive  bool    `json:"is_active"`
	Thumbnail string  `json:"thumbnail"`
}
