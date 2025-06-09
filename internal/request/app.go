package request

type DeleteManyRequest struct {
	IDs []string `json:"ids" binding:"required,min=1,dive,uuid"`
}