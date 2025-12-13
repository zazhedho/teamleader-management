package dto

type KPIItemCreate struct {
	PillarId          string   `json:"pillar_id" binding:"required"`
	Name              string   `json:"name" binding:"required,min=2,max=150"`
	Weight            float64  `json:"weight" binding:"required,gte=0,lte=100"`
	TargetValue       *float64 `json:"target_value" binding:"omitempty"`
	Unit              *string  `json:"unit" binding:"omitempty,max=50"`
	Frequency         *string  `json:"frequency" binding:"omitempty,oneof=DAILY WEEKLY MONTHLY QUARTERLY YEARLY"`
	InputSource       string   `json:"input_source" binding:"required,oneof=ADMIN_UPLOAD TL_INPUT SYSTEM"`
	AppliesToTL       bool     `json:"applies_to_tl"`
	AppliesToSalesman bool     `json:"applies_to_salesman"`
	Notes             *string  `json:"notes" binding:"omitempty,max=500"`
}

type KPIItemUpdate struct {
	PillarId          *string  `json:"pillar_id" binding:"omitempty"`
	Name              *string  `json:"name" binding:"omitempty,min=2,max=150"`
	Weight            *float64 `json:"weight" binding:"omitempty,gte=0,lte=100"`
	TargetValue       *float64 `json:"target_value" binding:"omitempty"`
	Unit              *string  `json:"unit" binding:"omitempty,max=50"`
	Frequency         *string  `json:"frequency" binding:"omitempty,oneof=DAILY WEEKLY MONTHLY QUARTERLY YEARLY"`
	InputSource       *string  `json:"input_source" binding:"omitempty,oneof=ADMIN_UPLOAD TL_INPUT SYSTEM"`
	AppliesToTL       *bool    `json:"applies_to_tl"`
	AppliesToSalesman *bool    `json:"applies_to_salesman"`
	Notes             *string  `json:"notes" binding:"omitempty,max=500"`
}

type PersonKPITargetUpsert struct {
	PersonId    string  `json:"person_id" binding:"required"`
	KPIItemId   string  `json:"kpi_item_id" binding:"required"`
	PeriodMonth int     `json:"period_month" binding:"required,min=1,max=12"`
	PeriodYear  int     `json:"period_year" binding:"required,min=2000"`
	TargetValue float64 `json:"target_value" binding:"required"`
}
