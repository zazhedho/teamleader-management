package domaintlsession

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func (TLSession) TableName() string {
	return "tl_sessions"
}

type TLSession struct {
	Id            string         `json:"id" gorm:"column:id;primaryKey"`
	PersonId      string         `json:"person_id" gorm:"column:person_id"`
	SessionType   string         `json:"session_type" gorm:"column:session_type"` // 'coaching', 'briefing'
	Date          time.Time      `json:"date" gorm:"column:date"`
	Notes         *string        `json:"notes,omitempty" gorm:"column:notes"`
	Attendees     pq.StringArray `json:"attendees,omitempty" gorm:"column:attendees;type:text[]"`
	DurationHours *float64       `json:"duration_hours,omitempty" gorm:"column:duration_hours"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
