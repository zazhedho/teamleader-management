package domaintlactivity

import (
	"time"

	"gorm.io/gorm"
)

func (TLDailyActivity) TableName() string {
	return "tl_daily_activities"
}

type TLDailyActivity struct {
	Id               string    `json:"id" gorm:"column:id;primaryKey"`
	PersonId         string    `json:"person_id" gorm:"column:person_id"`
	Date             time.Time `json:"date" gorm:"column:date"`
	ActivityType     string    `json:"activity_type" gorm:"column:activity_type"` // canvassing, pameran
	Kecamatan        *string   `json:"kecamatan,omitempty" gorm:"column:kecamatan"`
	Desa             *string   `json:"desa,omitempty" gorm:"column:desa"`
	GpsLat           *float64  `json:"gps_lat,omitempty" gorm:"column:gps_lat"`
	GpsLng           *float64  `json:"gps_lng,omitempty" gorm:"column:gps_lng"`
	DurationHours    *float64  `json:"duration_hours,omitempty" gorm:"column:duration_hours"` // for pameran
	ProspectCount    int       `json:"prospect_count" gorm:"column:prospect_count"`
	DealCount        int       `json:"deal_count" gorm:"column:deal_count"`
	MotorkuDownloads int       `json:"motorku_downloads" gorm:"column:motorku_downloads"`
	Notes            *string   `json:"notes,omitempty" gorm:"column:notes"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
