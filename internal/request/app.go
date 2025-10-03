package request

type DeleteManyRequest struct {
	IDs []int64 `json:"ids" binding:"required,min=1,dive"`
}
