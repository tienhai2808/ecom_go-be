package types

type UploadImageMessage struct {
	ImageID  int64     `json:"image_id"`
	FileName string    `json:"file_name"`
	FileData []byte `json:"file_data"`
}
