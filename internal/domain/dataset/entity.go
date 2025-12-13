package domaindataset

import (
	"time"

	"gorm.io/gorm"
)

func (DashboardDataset) TableName() string {
	return "dashboard_datasets"
}

type DashboardDataset struct {
	Id              string    `json:"id" gorm:"column:id;primaryKey"`
	Type            string    `json:"type" gorm:"column:type"`
	PeriodDate      time.Time `json:"period_date" gorm:"column:period_date"`
	PeriodMonth     int       `json:"period_month" gorm:"column:period_month"`
	PeriodYear      int       `json:"period_year" gorm:"column:period_year"`
	PeriodFrequency string    `json:"period_frequency" gorm:"column:period_frequency"`
	FileName        string    `json:"file_name" gorm:"column:file_name"`
	UploadedBy      string    `json:"uploaded_by" gorm:"column:uploaded_by"`
	UploadedAt      time.Time `json:"uploaded_at" gorm:"column:uploaded_at"`
	Status          string    `json:"status" gorm:"column:status"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
