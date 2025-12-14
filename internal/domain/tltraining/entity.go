package domaintltraining

import (
	"time"

	"gorm.io/gorm"
)

func (TLTrainingParticipation) TableName() string {
	return "tl_training_participations"
}

type TLTrainingParticipation struct {
	Id            string    `json:"id" gorm:"column:id;primaryKey"`
	TlPersonId    string    `json:"tl_person_id" gorm:"column:tl_person_id"`
	TrainingName  string    `json:"training_name" gorm:"column:training_name"`
	Date          time.Time `json:"date" gorm:"column:date"`
	SalesmanId    string    `json:"salesman_id" gorm:"column:salesman_id"`
	SalesmanName  string    `json:"salesman_name" gorm:"column:salesman_name"`
	Status        string    `json:"status" gorm:"column:status"` // hadir, tidak_hadir
	TrainingBatch string    `json:"training_batch" gorm:"column:training_batch;index"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	CreatedBy string         `json:"created_by" gorm:"column:created_by"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	UpdatedBy string         `json:"updated_by" gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	DeletedBy string         `json:"-"`
}
