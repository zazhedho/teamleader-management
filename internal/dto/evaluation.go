package dto

import "time"

// EvaluationCalculateRequest triggers evaluation calculation for a period
type EvaluationCalculateRequest struct {
	PeriodMonth int    `json:"period_month" binding:"required,min=1,max=12"`
	PeriodYear  int    `json:"period_year" binding:"required,min=2020"`
	PersonId    string `json:"person_id" binding:"omitempty,uuid4"` // Optional: calculate for specific TL, if empty calculate for all
}

// EvaluationResponse returns the evaluation result with breakdown
type EvaluationResponse struct {
	Id              string                 `json:"id"`
	PersonId        string                 `json:"person_id"`
	PersonName      string                 `json:"person_name,omitempty"`
	PeriodMonth     int                    `json:"period_month"`
	PeriodYear      int                    `json:"period_year"`
	TotalScore      float64                `json:"total_score"`
	PillarBreakdown []PillarScoreBreakdown `json:"pillar_breakdown,omitempty"`
	KpiBreakdown    []KpiScoreBreakdown    `json:"kpi_breakdown,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

// PillarScoreBreakdown shows scores aggregated by pillar
type PillarScoreBreakdown struct {
	PillarId       string  `json:"pillar_id"`
	PillarName     string  `json:"pillar_name"`
	PillarWeight   float64 `json:"pillar_weight"`    // e.g., 50.00 for Sales
	PillarScore    float64 `json:"pillar_score"`     // Weighted score from this pillar
	PillarMaxScore float64 `json:"pillar_max_score"` // Max possible score (same as weight)
}

// KpiScoreBreakdown shows detailed score for each KPI item
type KpiScoreBreakdown struct {
	KpiItemId        string   `json:"kpi_item_id"`
	KpiItemName      string   `json:"kpi_item_name"`
	PillarName       string   `json:"pillar_name"`
	Weight           float64  `json:"weight"`            // e.g., 25.00 for Quantity Activity
	ActualValue      *float64 `json:"actual_value"`      // Actual metric value
	TargetValue      *float64 `json:"target_value"`      // Target metric value (if applicable)
	AchievementRatio *float64 `json:"achievement_ratio"` // Actual/Target ratio (if applicable)
	Score            float64  `json:"score"`             // Weighted score from this KPI
	MaxScore         float64  `json:"max_score"`         // Max possible score (same as weight)
	Unit             *string  `json:"unit,omitempty"`    // e.g., "count", "percentage"
	InputSource      string   `json:"input_source"`      // "TL" or "ADMIN"
}

// EvaluationListResponse for paginated list
type EvaluationListResponse struct {
	Evaluations []EvaluationResponse `json:"evaluations"`
	Total       int64                `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

// LeaderboardEntry for ranking display
type LeaderboardEntry struct {
	Rank        int     `json:"rank"`
	PersonId    string  `json:"person_id"`
	PersonName  string  `json:"person_name"`
	DealerCode  string  `json:"dealer_code,omitempty"`
	TotalScore  float64 `json:"total_score"`
	PeriodMonth int     `json:"period_month"`
	PeriodYear  int     `json:"period_year"`
}

// LeaderboardResponse for ranking endpoints
type LeaderboardResponse struct {
	Period  string             `json:"period"` // e.g., "2025-12"
	Entries []LeaderboardEntry `json:"entries"`
	Total   int                `json:"total"`
}
