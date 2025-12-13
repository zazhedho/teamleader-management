package domainpillar

import (
	"time"

	"gorm.io/gorm"
)

func (Pillar) TableName() string {
	return "pillars"
}

type Pillar struct {
	Id          string  `json:"id" gorm:"column:id;primaryKey"`
	Name        string  `json:"name" gorm:"column:name"`
	Description *string `json:"description,omitempty" gorm:"column:description"`
	Weight      float64 `json:"weight" gorm:"column:weight"` // percent, supports decimals

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty" gorm:"column:updated_at"`
	UpdatedBy *string        `json:"updated_by,omitempty" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
