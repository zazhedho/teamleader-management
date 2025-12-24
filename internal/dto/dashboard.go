package dto

import "time"

// ========================================
// TL DASHBOARD
// ========================================

// TLDashboardResponse provides overview for a TL
type TLDashboardResponse struct {
	PersonInfo        PersonInfo             `json:"person_info"`
	CurrentEvaluation *EvaluationSummary     `json:"current_evaluation"`
	RecentActivities  []RecentActivityItem   `json:"recent_activities"`
	PerformanceTrend  []PerformanceTrendItem `json:"performance_trend"`
	QuickStats        TLQuickStats           `json:"quick_stats"`
	Ranking           *RankingInfo           `json:"ranking,omitempty"`
}

// PersonInfo contains TL basic information
type PersonInfo struct {
	PersonId   string  `json:"person_id"`
	Name       string  `json:"name"`
	HondaId    string  `json:"honda_id"`
	DealerCode *string `json:"dealer_code,omitempty"`
	Role       string  `json:"role"`
}

// EvaluationSummary is a simplified evaluation for dashboard
type EvaluationSummary struct {
	EvaluationId string               `json:"evaluation_id"`
	PeriodMonth  int                  `json:"period_month"`
	PeriodYear   int                  `json:"period_year"`
	TotalScore   float64              `json:"total_score"`
	PillarScores []PillarScoreSummary `json:"pillar_scores"`
	TopKPIs      []KPISummary         `json:"top_kpis"`  // Top 3 KPIs
	WeakKPIs     []KPISummary         `json:"weak_kpis"` // Bottom 3 KPIs
	LastUpdated  time.Time            `json:"last_updated"`
}

// PillarScoreSummary simplified pillar score
type PillarScoreSummary struct {
	Name       string  `json:"name"`
	Score      float64 `json:"score"`
	MaxScore   float64 `json:"max_score"`
	Percentage float64 `json:"percentage"` // score/max_score * 100
}

// KPISummary simplified KPI info
type KPISummary struct {
	Name        string   `json:"name"`
	Score       float64  `json:"score"`
	MaxScore    float64  `json:"max_score"`
	ActualValue *float64 `json:"actual_value,omitempty"`
	TargetValue *float64 `json:"target_value,omitempty"`
	Percentage  float64  `json:"percentage"`
}

// RecentActivityItem shows recent TL activities
type RecentActivityItem struct {
	Type        string    `json:"type"` // "activity", "attendance", "coaching", "briefing", "training"
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Details     string    `json:"details,omitempty"`
}

// PerformanceTrendItem shows score over time
type PerformanceTrendItem struct {
	PeriodMonth int     `json:"period_month"`
	PeriodYear  int     `json:"period_year"`
	PeriodLabel string  `json:"period_label"` // e.g., "Dec 2025"
	TotalScore  float64 `json:"total_score"`
}

// TLQuickStats quick statistics for current month
type TLQuickStats struct {
	ActivitiesCount int     `json:"activities_count"`
	CoachingCount   int     `json:"coaching_count"`
	BriefingCount   int     `json:"briefing_count"`
	TrainingCount   int     `json:"training_count"`
	TeamSize        int     `json:"team_size"`
	AttendanceRate  float64 `json:"attendance_rate"` // percentage
}

// RankingInfo shows TL's rank among peers
type RankingInfo struct {
	Rank       int `json:"rank"`
	TotalTLs   int `json:"total_tls"`
	Percentile int `json:"percentile"` // Top X%
}

// ========================================
// ADMIN ANALYTICS
// ========================================

// AdminAnalyticsResponse provides aggregated insights
type AdminAnalyticsResponse struct {
	Period           string             `json:"period"` // e.g., "2025-12"
	OverallStats     OverallStatistics  `json:"overall_stats"`
	TopPerformers    []LeaderboardEntry `json:"top_performers"`
	BottomPerformers []LeaderboardEntry `json:"bottom_performers"`
	PillarAnalysis   []PillarAnalysis   `json:"pillar_analysis"`
	TrendComparison  []PeriodComparison `json:"trend_comparison"`
	Distributions    ScoreDistribution  `json:"distributions"`
}

// OverallStatistics aggregated stats for all TLs
type OverallStatistics struct {
	TotalTLs          int     `json:"total_tls"`
	AverageScore      float64 `json:"average_score"`
	MedianScore       float64 `json:"median_score"`
	HighestScore      float64 `json:"highest_score"`
	LowestScore       float64 `json:"lowest_score"`
	StandardDeviation float64 `json:"standard_deviation"`
}

// PillarAnalysis shows performance by pillar
type PillarAnalysis struct {
	PillarName     string  `json:"pillar_name"`
	AverageScore   float64 `json:"average_score"`
	MaxPossible    float64 `json:"max_possible"`
	AchievementPct float64 `json:"achievement_pct"` // average/max * 100
	TopScorer      string  `json:"top_scorer"`      // TL name
	TopScore       float64 `json:"top_score"`
}

// PeriodComparison compares metrics across periods
type PeriodComparison struct {
	PeriodMonth      int      `json:"period_month"`
	PeriodYear       int      `json:"period_year"`
	PeriodLabel      string   `json:"period_label"`
	AverageScore     float64  `json:"average_score"`
	ParticipatingTLs int      `json:"participating_tls"`
	Change           *float64 `json:"change,omitempty"` // vs previous period
}

// ScoreDistribution shows how scores are distributed
type ScoreDistribution struct {
	Range0_20   int `json:"range_0_20"`   // 0-20 points
	Range20_40  int `json:"range_20_40"`  // 20-40 points
	Range40_60  int `json:"range_40_60"`  // 40-60 points
	Range60_80  int `json:"range_60_80"`  // 60-80 points
	Range80_100 int `json:"range_80_100"` // 80-100 points
}

// ========================================
// PERFORMANCE COMPARISON
// ========================================

// TeamComparisonRequest for comparing multiple TLs
type TeamComparisonRequest struct {
	PersonIds   []string `json:"person_ids" binding:"required,min=2,max=10"`
	PeriodMonth int      `json:"period_month" binding:"required,min=1,max=12"`
	PeriodYear  int      `json:"period_year" binding:"required,min=2020"`
}

// TeamComparisonResponse compares multiple TLs side by side
type TeamComparisonResponse struct {
	Period      string                  `json:"period"`
	Comparisons []TLComparisonItem      `json:"comparisons"`
	PillarChart []PillarComparisonChart `json:"pillar_chart"`
}

// TLComparisonItem individual TL in comparison
type TLComparisonItem struct {
	PersonId     string               `json:"person_id"`
	PersonName   string               `json:"person_name"`
	TotalScore   float64              `json:"total_score"`
	Rank         int                  `json:"rank"`
	PillarScores []PillarScoreSummary `json:"pillar_scores"`
}

// PillarComparisonChart data for chart visualization
type PillarComparisonChart struct {
	PillarName string          `json:"pillar_name"`
	TLScores   []TLPillarScore `json:"tl_scores"`
}

// TLPillarScore individual score in chart
type TLPillarScore struct {
	PersonName string  `json:"person_name"`
	Score      float64 `json:"score"`
}
