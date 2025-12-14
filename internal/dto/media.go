package dto

// Media DTOs
type MediaCreate struct {
	EntityType string  `json:"entity_type" binding:"required"`
	EntityId   string  `json:"entity_id" binding:"required"`
	FileUrl    string  `json:"file_url" binding:"required,url"`
	FileName   string  `json:"file_name" binding:"required,min=1,max=255"`
	FileType   *string `json:"file_type" binding:"omitempty"`
	FileSize   *int64  `json:"file_size" binding:"omitempty,min=0"`
}

type MediaAttach struct {
	FileUrl  string  `json:"file_url" binding:"required,url"`
	FileName string  `json:"file_name" binding:"required,min=1,max=255"`
	FileType *string `json:"file_type" binding:"omitempty"`
	FileSize *int64  `json:"file_size" binding:"omitempty,min=0"`
}
