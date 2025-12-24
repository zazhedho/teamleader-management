package domainevaluation

import "time"

// Evaluation represents the overall evaluation result for a person in a specific period
type Evaluation struct {
	Id                 string    `json:"id" gorm:"column:id;primaryKey"`
	EvaluationPeriodId string    `json:"evaluation_period_id" gorm:"column:evaluation_period_id"`
	PersonId           string    `json:"person_id" gorm:"column:person_id"`
	TotalScore         float64   `json:"total_score" gorm:"column:total_score"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at"`

	// Relations (not stored in DB, loaded via joins)
	Period  *EvaluationPeriod  `json:"period,omitempty" gorm:"foreignKey:EvaluationPeriodId"`
	Details []EvaluationDetail `json:"details,omitempty" gorm:"foreignKey:EvaluationId"`
}

func (Evaluation) TableName() string {
	return "evaluations"
}

// EvaluationDetail represents the detailed score for each KPI item
type EvaluationDetail struct {
	Id               string   `json:"id" gorm:"column:id;primaryKey"`
	EvaluationId     string   `json:"evaluation_id" gorm:"column:evaluation_id"`
	KpiItemId        string   `json:"kpi_item_id" gorm:"column:kpi_item_id"`
	ActualValue      *float64 `json:"actual_value" gorm:"column:actual_value"`
	AchievementRatio *float64 `json:"achievement_ratio" gorm:"column:achievement_ratio"`
	Score            float64  `json:"score" gorm:"column:score"`
}

func (EvaluationDetail) TableName() string {
	return "evaluation_details"
}

// EvaluationPeriod represents the evaluation period (monthly)
type EvaluationPeriod struct {
	Id          string    `json:"id" gorm:"column:id;primaryKey"`
	PeriodMonth int       `json:"period_month" gorm:"column:period_month"`
	PeriodYear  int       `json:"period_year" gorm:"column:period_year"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
}

func (EvaluationPeriod) TableName() string {
	return "evaluation_periods"
}
