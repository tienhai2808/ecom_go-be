package types

type UploadImageMessage struct {
	ImageID  int64  `json:"image_id"`
	FileName string `json:"file_name"`
	FileData []byte `json:"file_data"`
}

type ProductSearchResult struct {
	IDs        []int64 `json:"product_ids"`
	Total      int64   `json:"total"`
	Page       uint32  `json:"page"`
	Limit      uint32  `json:"limit"`
	TotalPages int64   `json:"total_pages"`
	HasPrev    bool    `json:"has_prev"`
	HasNext    bool    `json:"has_next"`
}
