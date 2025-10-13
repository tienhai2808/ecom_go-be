package response

type ApiResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

type MetaResponse struct {
	Total      int64  `json:"total"`
	Page       uint32 `json:"page"`
	Limit      uint32 `json:"limit"`
	TotalPages int64  `json:"total_pages"`
	HasPrev    bool   `json:"has_prev"`
	HasNext    bool   `json:"has_next"`
}
