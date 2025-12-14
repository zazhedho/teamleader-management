package domainmedia

import (
	"time"

	"gorm.io/gorm"
)

func (Media) TableName() string {
	return "media"
}

// Media is a generic table for storing file attachments across the system
// Can be used by TL activities, datasets, user profiles, etc.
type Media struct {
	Id           string  `json:"id" gorm:"column:id;primaryKey"`
	EntityType   string  `json:"entity_type" gorm:"column:entity_type;index"` // e.g., 'tl_activity', 'tl_session', 'dataset', 'user_profile'
	EntityId     string  `json:"entity_id" gorm:"column:entity_id;index"`     // FK to related entity
	FileUrl      string  `json:"file_url" gorm:"column:file_url"`
	FileName     string  `json:"file_name" gorm:"column:file_name"`
	FileType     *string `json:"file_type,omitempty" gorm:"column:file_type"` // e.g., 'image/jpeg', 'application/pdf'
	FileSize     *int64  `json:"file_size,omitempty" gorm:"column:file_size"` // in bytes
	DisplayOrder int     `json:"display_order" gorm:"column:display_order;default:1"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
