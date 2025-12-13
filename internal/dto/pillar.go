package dto

type PillarCreate struct {
	Name        string  `json:"name" binding:"required,min=2,max=150"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Weight      float64 `json:"weight" binding:"required,gte=0,lte=100"`
}

type PillarUpdate struct {
	Name        *string  `json:"name" binding:"omitempty,min=2,max=150"`
	Description *string  `json:"description" binding:"omitempty,max=500"`
	Weight      *float64 `json:"weight" binding:"omitempty,gte=0,lte=100"`
}
