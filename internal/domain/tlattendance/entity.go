package domaintlattendance

import (
	"time"

	"gorm.io/gorm"
)

func (TLAttendanceRecord) TableName() string {
	return "tl_attendance_records"
}

type TLAttendanceRecord struct {
	Id             string    `json:"id" gorm:"column:id;primaryKey"`
	TlPersonId     string    `json:"tl_person_id" gorm:"column:tl_person_id"`
	SalesmanId     string    `json:"salesman_id" gorm:"column:salesman_id"`
	SalesmanName   string    `json:"salesman_name" gorm:"column:salesman_name"`
	Date           time.Time `json:"date" gorm:"column:date"`
	Status         string    `json:"status" gorm:"column:status"` // hadir, tidak_hadir
	RecordUniqueId string    `json:"record_unique_id" gorm:"column:record_unique_id;index"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
