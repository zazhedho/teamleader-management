package domainperson

import (
	"time"

	"gorm.io/gorm"
)

func (Person) TableName() string {
	return "persons"
}

type Person struct {
	Id         string  `json:"id" gorm:"column:id;primaryKey"`
	HondaId    string  `json:"honda_id" gorm:"column:honda_id"`
	Name       string  `json:"name" gorm:"column:name"`
	JobTitle   *string `json:"job_title,omitempty" gorm:"column:job_title"`
	Role       string  `json:"role" gorm:"column:role"`
	DealerCode *string `json:"dealer_code,omitempty" gorm:"column:dealer_code"`
	Active     bool    `json:"active" gorm:"column:active"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
