package domainkpiitem

import (
	"time"

	"gorm.io/gorm"
)

func (KPIItem) TableName() string {
	return "kpi_items"
}

type KPIItem struct {
	Id                string   `json:"id" gorm:"column:id;primaryKey"`
	PillarId          string   `json:"pillar_id" gorm:"column:pillar_id"`
	Name              string   `json:"name" gorm:"column:name"`
	Weight            float64  `json:"weight" gorm:"column:weight"` // percent 0-100
	TargetValue       *float64 `json:"target_value,omitempty" gorm:"column:target_value"`
	Unit              *string  `json:"unit,omitempty" gorm:"column:unit"`           // x/hari | point | ratio | percentage | etc...
	Frequency         *string  `json:"frequency,omitempty" gorm:"column:frequency"` // DAILY | WEEKLY | MONTHLY | YEARLY
	InputSource       string   `json:"input_source" gorm:"column:input_source"`     // ADMIN | TL | SYSTEM
	AppliesToTL       bool     `json:"applies_to_tl" gorm:"column:applies_to_tl"`
	AppliesToSalesman bool     `json:"applies_to_salesman" gorm:"column:applies_to_salesman"`
	Notes             *string  `json:"notes,omitempty" gorm:"column:notes"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}

func (PersonKPITarget) TableName() string {
	return "person_kpi_targets"
}

type PersonKPITarget struct {
	Id          string  `json:"id" gorm:"column:id;primaryKey"`
	PersonId    string  `json:"person_id" gorm:"column:person_id"`
	KPIItemId   string  `json:"kpi_item_id" gorm:"column:kpi_item_id"`
	PeriodMonth int     `json:"period_month" gorm:"column:period_month"`
	PeriodYear  int     `json:"period_year" gorm:"column:period_year"`
	TargetValue float64 `json:"target_value" gorm:"column:target_value"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
