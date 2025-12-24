package serviceevaluation

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MetricAggregator aggregates metrics from various sources for evaluation
type MetricAggregator struct {
	DB *gorm.DB
}

func NewMetricAggregator(db *gorm.DB) *MetricAggregator {
	return &MetricAggregator{DB: db}
}

// MetricValue represents a raw metric value with metadata
type MetricValue struct {
	Value       float64
	Unit        string
	Source      string
	Description string
}

// GetMetricsForPerson retrieves all metrics for a person in a specific period
func (m *MetricAggregator) GetMetricsForPerson(personId string, periodMonth int, periodYear int) (map[string]*MetricValue, error) {
	metrics := make(map[string]*MetricValue)

	// Calculate period boundaries
	startDate := time.Date(periodYear, time.Month(periodMonth), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// ========================================
	// SALES PERFORMANCE (50%)
	// ========================================

	// 1. Quantity Activity (25%) - Count of daily activities
	activityCount, err := m.countDailyActivities(personId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity count: %w", err)
	}
	metrics["quantity_activity"] = &MetricValue{
		Value:       float64(activityCount),
		Unit:        "count",
		Source:      "TL",
		Description: "Number of promotional activities (canvassing + pameran)",
	}

	// 2. Sales FLP (25%) - From admin dataset
	salesFLP, err := m.getSalesFLP(personId, periodMonth, periodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales FLP: %w", err)
	}
	metrics["sales_flp"] = &MetricValue{
		Value:       salesFLP,
		Unit:        "amount",
		Source:      "ADMIN",
		Description: "Sales FLP amount",
	}

	// ========================================
	// LEADERSHIP (15%)
	// ========================================

	// 3. Discipline & Attendance (2.5%) - Attendance percentage
	attendancePercentage, err := m.getAttendancePercentage(personId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendance: %w", err)
	}
	metrics["attendance"] = &MetricValue{
		Value:       attendancePercentage,
		Unit:        "percentage",
		Source:      "TL",
		Description: "Team attendance rate",
	}

	// 4. Coaching Sessions (2.5%) - Count of coaching sessions
	coachingCount, err := m.countSessions(personId, startDate, endDate, "coaching")
	if err != nil {
		return nil, fmt.Errorf("failed to get coaching count: %w", err)
	}
	metrics["coaching_sessions"] = &MetricValue{
		Value:       float64(coachingCount),
		Unit:        "count",
		Source:      "TL",
		Description: "Number of coaching sessions",
	}

	// 5. Briefing Sessions (2.5%) - Count of briefing sessions
	briefingCount, err := m.countSessions(personId, startDate, endDate, "briefing")
	if err != nil {
		return nil, fmt.Errorf("failed to get briefing count: %w", err)
	}
	metrics["briefing_sessions"] = &MetricValue{
		Value:       float64(briefingCount),
		Unit:        "count",
		Source:      "TL",
		Description: "Number of briefing sessions",
	}

	// 6. Team Size (7.5%) - Number of team members
	teamSize, err := m.getTeamSize(personId)
	if err != nil {
		return nil, fmt.Errorf("failed to get team size: %w", err)
	}
	metrics["team_size"] = &MetricValue{
		Value:       float64(teamSize),
		Unit:        "count",
		Source:      "TL",
		Description: "Number of team members",
	}

	// ========================================
	// DEVELOPMENT (10%)
	// ========================================

	// 7. Quiz Score (5%) - From admin dataset
	quizScore, err := m.getQuizScore(personId, periodMonth, periodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz score: %w", err)
	}
	metrics["quiz_score"] = &MetricValue{
		Value:       quizScore,
		Unit:        "score",
		Source:      "ADMIN",
		Description: "Quiz result score",
	}

	// 8. Training Participation (5%) - Count of training participations
	trainingCount, err := m.countTrainingParticipations(personId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get training count: %w", err)
	}
	metrics["training_participation"] = &MetricValue{
		Value:       float64(trainingCount),
		Unit:        "count",
		Source:      "TL",
		Description: "Number of training sessions attended",
	}

	// ========================================
	// DIGITALIZATION (25%)
	// ========================================

	// 9. Login Apple (5%) - From admin dataset
	appleLogins, err := m.getAppleLogins(personId, periodMonth, periodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get apple logins: %w", err)
	}
	metrics["apple_logins"] = &MetricValue{
		Value:       float64(appleLogins),
		Unit:        "count",
		Source:      "ADMIN",
		Description: "Number of Apple app logins",
	}

	// 10. Point Apple (5%) - From admin dataset
	applePoints, err := m.getApplePoints(personId, periodMonth, periodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get apple points: %w", err)
	}
	metrics["apple_points"] = &MetricValue{
		Value:       float64(applePoints),
		Unit:        "points",
		Source:      "ADMIN",
		Description: "Apple app points earned",
	}

	// 11. Point My Hero (5%) - From admin dataset
	myHeroPoints, err := m.getMyHeroPoints(personId, periodMonth, periodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get my hero points: %w", err)
	}
	metrics["myhero_points"] = &MetricValue{
		Value:       float64(myHeroPoints),
		Unit:        "points",
		Source:      "ADMIN",
		Description: "My Hero app points earned",
	}

	// 12. Total Prospects (5%) - From admin dataset
	totalProspects, err := m.getTotalProspects(personId, periodMonth, periodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get total prospects: %w", err)
	}
	metrics["total_prospects"] = &MetricValue{
		Value:       float64(totalProspects),
		Unit:        "count",
		Source:      "ADMIN",
		Description: "Total number of prospects",
	}

	// 13. Prospect Ratio (5%) - Calculated (requires target)
	// Note: This will be calculated in the calculator service based on target
	metrics["prospect_ratio"] = &MetricValue{
		Value:       0, // Will be calculated later
		Unit:        "ratio",
		Source:      "ADMIN",
		Description: "Prospect achievement ratio",
	}

	return metrics, nil
}

// ========================================
// HELPER FUNCTIONS - Query each data source
// ========================================

func (m *MetricAggregator) countDailyActivities(personId string, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := m.DB.Table("tl_daily_activities").
		Where("person_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", personId, startDate, endDate).
		Count(&count).Error
	return count, err
}

func (m *MetricAggregator) getSalesFLP(personId string, periodMonth int, periodYear int) (float64, error) {
	var result struct {
		FLPAmount float64
	}

	err := m.DB.Table("sales_flp sf").
		Select("COALESCE(SUM(sf.flp_amount), 0) as flp_amount").
		Joins("INNER JOIN dashboard_datasets dd ON sf.dataset_id = dd.id").
		Where("sf.person_id = ? AND dd.period_month = ? AND dd.period_year = ?", personId, periodMonth, periodYear).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.FLPAmount, nil
}

func (m *MetricAggregator) getAttendancePercentage(personId string, startDate, endDate time.Time) (float64, error) {
	var result struct {
		AvgPercentage float64
	}

	// Calculate average attendance percentage across all records in period
	err := m.DB.Table("tl_attendance_records").
		Select("COALESCE(AVG(CASE WHEN status = 'hadir' THEN 100.0 ELSE 0.0 END), 0) as avg_percentage").
		Where("tl_person_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", personId, startDate, endDate).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.AvgPercentage, nil
}

func (m *MetricAggregator) countSessions(personId string, startDate, endDate time.Time, sessionType string) (int64, error) {
	var count int64
	err := m.DB.Table("tl_sessions").
		Where("person_id = ? AND session_type = ? AND date >= ? AND date <= ? AND deleted_at IS NULL",
			personId, sessionType, startDate, endDate).
		Count(&count).Error
	return count, err
}

func (m *MetricAggregator) getTeamSize(personId string) (int64, error) {
	// Count persons who have this TL as their supervisor
	// Note: This assumes there's a supervisor_id or similar field
	// Adjust the query based on actual schema
	var count int64
	err := m.DB.Table("persons").
		Where("role = 'salesman' AND active = true").
		// TODO: Add supervisor_id filter when schema is updated
		Count(&count).Error
	return count, err
}

func (m *MetricAggregator) getQuizScore(personId string, periodMonth int, periodYear int) (float64, error) {
	var result struct {
		Score float64
	}

	err := m.DB.Table("quiz_results qr").
		Select("COALESCE(AVG(qr.score), 0) as score").
		Joins("INNER JOIN dashboard_datasets dd ON qr.dataset_id = dd.id").
		Where("qr.person_id = ? AND dd.period_month = ? AND dd.period_year = ?", personId, periodMonth, periodYear).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.Score, nil
}

func (m *MetricAggregator) countTrainingParticipations(personId string, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := m.DB.Table("tl_training_participations").
		Where("person_id = ? AND start_date >= ? AND start_date <= ? AND deleted_at IS NULL",
			personId, startDate, endDate).
		Count(&count).Error
	return count, err
}

func (m *MetricAggregator) getAppleLogins(personId string, periodMonth int, periodYear int) (int, error) {
	var result struct {
		LoginCount int
	}

	err := m.DB.Table("apple_logins al").
		Select("COALESCE(SUM(al.login_count), 0) as login_count").
		Joins("INNER JOIN dashboard_datasets dd ON al.dataset_id = dd.id").
		Where("al.person_id = ? AND dd.period_month = ? AND dd.period_year = ?", personId, periodMonth, periodYear).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.LoginCount, nil
}

func (m *MetricAggregator) getApplePoints(personId string, periodMonth int, periodYear int) (int, error) {
	var result struct {
		Points int
	}

	err := m.DB.Table("apple_points ap").
		Select("COALESCE(SUM(ap.points), 0) as points").
		Joins("INNER JOIN dashboard_datasets dd ON ap.dataset_id = dd.id").
		Where("ap.person_id = ? AND dd.period_month = ? AND dd.period_year = ?", personId, periodMonth, periodYear).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.Points, nil
}

func (m *MetricAggregator) getMyHeroPoints(personId string, periodMonth int, periodYear int) (int, error) {
	var result struct {
		Points int
	}

	err := m.DB.Table("myhero_points mp").
		Select("COALESCE(SUM(mp.points), 0) as points").
		Joins("INNER JOIN dashboard_datasets dd ON mp.dataset_id = dd.id").
		Where("mp.person_id = ? AND dd.period_month = ? AND dd.period_year = ?", personId, periodMonth, periodYear).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.Points, nil
}

func (m *MetricAggregator) getTotalProspects(personId string, periodMonth int, periodYear int) (int, error) {
	var result struct {
		ProspectCount int
	}

	err := m.DB.Table("prospects p").
		Select("COALESCE(SUM(p.prospect_count), 0) as prospect_count").
		Joins("INNER JOIN dashboard_datasets dd ON p.dataset_id = dd.id").
		Where("p.person_id = ? AND dd.period_month = ? AND dd.period_year = ?", personId, periodMonth, periodYear).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}
	return result.ProspectCount, nil
}
