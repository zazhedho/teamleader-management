package domainmetric

import (
	"time"

	"gorm.io/gorm"
)

type QuizResult struct {
	Id         string   `json:"id" gorm:"column:id;primaryKey"`
	DatasetId  string   `json:"dataset_id" gorm:"column:dataset_id"`
	PersonId   string   `json:"person_id" gorm:"column:person_id"`
	HondaId    string   `json:"honda_id" gorm:"column:honda_id"`
	DealerCode *string  `json:"dealer_code,omitempty" gorm:"column:dealer_code"`
	Score      *float64 `json:"score,omitempty" gorm:"column:score"`
	PassStatus *string  `json:"pass_status,omitempty" gorm:"column:pass_status"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}

type AppleLogin struct {
	Id          string    `json:"id" gorm:"column:id;primaryKey"`
	DatasetId   string    `json:"dataset_id" gorm:"column:dataset_id"`
	PersonId    string    `json:"person_id" gorm:"column:person_id"`
	HondaId     string    `json:"honda_id" gorm:"column:honda_id"`
	DealerCode  *string   `json:"dealer_code,omitempty" gorm:"column:dealer_code"`
	LoginDate   time.Time `json:"login_date" gorm:"column:login_date"`
	MorningDone bool      `json:"morning_done" gorm:"column:morning_done"`
	EveningDone bool      `json:"evening_done" gorm:"column:evening_done"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}

type SalesFLP struct {
	Id         string  `json:"id" gorm:"column:id;primaryKey"`
	DatasetId  string  `json:"dataset_id" gorm:"column:dataset_id"`
	PersonId   string  `json:"person_id" gorm:"column:person_id"`
	HondaId    string  `json:"honda_id" gorm:"column:honda_id"`
	DealerCode *string `json:"dealer_code,omitempty" gorm:"column:dealer_code"`
	Amount     float64 `json:"amount" gorm:"column:flp_amount"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}

type ApplePoint struct {
	Id         string  `json:"id" gorm:"column:id;primaryKey"`
	DatasetId  string  `json:"dataset_id" gorm:"column:dataset_id"`
	PersonId   string  `json:"person_id" gorm:"column:person_id"`
	HondaId    string  `json:"honda_id" gorm:"column:honda_id"`
	DealerCode *string `json:"dealer_code,omitempty" gorm:"column:dealer_code"`
	Points     int     `json:"points" gorm:"column:points"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}

type MyHeroPoint struct {
	Id         string  `json:"id" gorm:"column:id;primaryKey"`
	DatasetId  string  `json:"dataset_id" gorm:"column:dataset_id"`
	PersonId   string  `json:"person_id" gorm:"column:person_id"`
	HondaId    string  `json:"honda_id" gorm:"column:honda_id"`
	DealerCode *string `json:"dealer_code,omitempty" gorm:"column:dealer_code"`
	Points     int     `json:"points" gorm:"column:points"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}

type Prospect struct {
	Id            string `json:"id" gorm:"column:id;primaryKey"`
	DatasetId     string `json:"dataset_id" gorm:"column:dataset_id"`
	PersonId      string `json:"person_id" gorm:"column:person_id"`
	HondaId       string `json:"honda_id" gorm:"column:honda_id"`
	ProspectCount int    `json:"prospect_count" gorm:"column:prospect_count"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
