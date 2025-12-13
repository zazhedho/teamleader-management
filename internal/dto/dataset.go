package dto

type DatasetUploadRequest struct {
	PeriodDate      string `form:"period_date" binding:"required"`                // YYYY-MM-DD anchor date
	PeriodMonth     int    `form:"period_month" binding:"omitempty,min=1,max=12"` // optional derived
	PeriodYear      int    `form:"period_year" binding:"omitempty,min=2000"`      // optional derived
	PeriodFrequency string `form:"period_frequency" binding:"omitempty,oneof=DAILY WEEKLY MONTHLY QUARTERLY YEARLY"`
	Type            string `form:"type" binding:"required"`
}

type DatasetStatusUpdate struct {
	Status string `json:"status" binding:"required"`
}
