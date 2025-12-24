package serviceevaluation

import (
	"fmt"

	domainevaluation "teamleader-management/internal/domain/evaluation"
	domainkpiitem "teamleader-management/internal/domain/kpiitem"
	"teamleader-management/utils"

	"gorm.io/gorm"
)

// EvaluationCalculator calculates evaluation scores based on metrics and KPI weights
type EvaluationCalculator struct {
	DB               *gorm.DB
	MetricAggregator *MetricAggregator
}

func NewEvaluationCalculator(db *gorm.DB) *EvaluationCalculator {
	return &EvaluationCalculator{
		DB:               db,
		MetricAggregator: NewMetricAggregator(db),
	}
}

// CalculationResult contains the full evaluation result with details
type CalculationResult struct {
	TotalScore float64
	Details    []domainevaluation.EvaluationDetail
}

// Calculate performs the evaluation calculation for a person in a period
func (c *EvaluationCalculator) Calculate(personId string, periodMonth int, periodYear int) (*CalculationResult, error) {
	// 1. Get all KPI items for TL
	kpiItems, err := c.getKPIItems()
	if err != nil {
		return nil, fmt.Errorf("failed to get KPI items: %w", err)
	}

	if len(kpiItems) == 0 {
		return nil, fmt.Errorf("no KPI items found for evaluation")
	}

	// 2. Get metrics for the person
	metrics, err := c.MetricAggregator.GetMetricsForPerson(personId, periodMonth, periodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate metrics: %w", err)
	}

	// 3. Calculate score for each KPI
	var details []domainevaluation.EvaluationDetail
	var totalScore float64

	for _, kpi := range kpiItems {
		detail, score := c.calculateKPIScore(kpi, metrics)
		details = append(details, detail)
		totalScore += score
	}

	return &CalculationResult{
		TotalScore: totalScore,
		Details:    details,
	}, nil
}

// calculateKPIScore calculates score for a single KPI item
func (c *EvaluationCalculator) calculateKPIScore(kpi domainkpiitem.KPIItem, metrics map[string]*MetricValue) (domainevaluation.EvaluationDetail, float64) {
	detail := domainevaluation.EvaluationDetail{
		Id:        utils.CreateUUID(),
		KpiItemId: kpi.Id,
		Score:     0,
	}

	// Map KPI names to metric keys
	metricKey := c.getMetricKey(kpi.Name)
	metric, exists := metrics[metricKey]

	if !exists {
		// No data for this KPI, score is 0
		return detail, 0
	}

	actualValue := metric.Value
	detail.ActualValue = &actualValue

	// Calculate score based on KPI type
	score := c.calculateScore(kpi, actualValue)

	// Calculate achievement ratio if target is set
	if kpi.TargetValue != nil && *kpi.TargetValue > 0 {
		achievementRatio := actualValue / *kpi.TargetValue * 100
		detail.AchievementRatio = &achievementRatio
	}

	detail.Score = score
	return detail, score
}

// calculateScore determines the score based on achievement
func (c *EvaluationCalculator) calculateScore(kpi domainkpiitem.KPIItem, actualValue float64) float64 {
	weight := kpi.Weight

	// If no target is set, use simple presence scoring
	if kpi.TargetValue == nil || *kpi.TargetValue == 0 {
		// For non-target metrics, give full score if value > 0
		if actualValue > 0 {
			return weight
		}
		return 0
	}

	// Calculate achievement ratio
	achievementRatio := actualValue / *kpi.TargetValue

	// Cap at 100% (don't give bonus for exceeding target)
	if achievementRatio > 1.0 {
		achievementRatio = 1.0
	}

	// Score is proportional to achievement
	score := weight * achievementRatio

	return score
}

// getMetricKey maps KPI name to metric key
func (c *EvaluationCalculator) getMetricKey(kpiName string) string {
	// Map KPI item names to metric keys
	mapping := map[string]string{
		"Quantity Activity":        "quantity_activity",
		"Sales FLP":                "sales_flp",
		"Disiplin & Kehadiran Tim": "attendance",
		"Sesi Coaching":            "coaching_sessions",
		"Sesi Briefing":            "briefing_sessions",
		"Jumlah Tim":               "team_size",
		"Kuis":                     "quiz_score",
		"Partisipasi Training":     "training_participation",
		"Login Apple":              "apple_logins",
		"Point Apple":              "apple_points",
		"Point my Hero":            "myhero_points",
		"Jumlah Prospek":           "total_prospects",
		"Ratio Prospek":            "prospect_ratio",
	}

	if key, found := mapping[kpiName]; found {
		return key
	}

	// Fallback: use lowercase with underscores
	return kpiName
}

// getKPIItems retrieves all KPI items for TL evaluation
func (c *EvaluationCalculator) getKPIItems() ([]domainkpiitem.KPIItem, error) {
	var kpiItems []domainkpiitem.KPIItem

	err := c.DB.Where("applies_to_tl = ?", true).
		Order("pillar_id, name").
		Find(&kpiItems).Error

	if err != nil {
		return nil, err
	}

	return kpiItems, nil
}
