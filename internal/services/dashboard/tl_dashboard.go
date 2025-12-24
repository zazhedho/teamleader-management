package servicedashboard

import (
	"fmt"
	"sort"
	"time"

	domainevaluation "teamleader-management/internal/domain/evaluation"
	domainperson "teamleader-management/internal/domain/person"
	"teamleader-management/internal/dto"
	interfaceevaluation "teamleader-management/internal/interfaces/evaluation"

	"gorm.io/gorm"
)

type TLDashboardService struct {
	DB       *gorm.DB
	EvalRepo interfaceevaluation.RepoEvaluationInterface
}

func NewTLDashboardService(db *gorm.DB, evalRepo interfaceevaluation.RepoEvaluationInterface) *TLDashboardService {
	return &TLDashboardService{
		DB:       db,
		EvalRepo: evalRepo,
	}
}

// GetDashboard retrieves complete dashboard for a TL
func (s *TLDashboardService) GetDashboard(personId string, periodMonth int, periodYear int) (dto.TLDashboardResponse, error) {
	// Get person info
	personInfo, err := s.getPersonInfo(personId)
	if err != nil {
		return dto.TLDashboardResponse{}, fmt.Errorf("failed to get person info: %w", err)
	}

	// Get current evaluation
	period, err := s.EvalRepo.GetPeriodByMonthYear(periodMonth, periodYear)
	var currentEval *dto.EvaluationSummary
	if err == nil {
		eval, err := s.EvalRepo.GetByPersonAndPeriod(personId, period.Id)
		if err == nil {
			summary, _ := s.buildEvaluationSummary(eval)
			currentEval = &summary
		}
	}

	// Get recent activities
	recentActivities := s.getRecentActivities(personId, 10)

	// Get performance trend (last 6 months)
	performanceTrend := s.getPerformanceTrend(personId, 6)

	// Get quick stats for current month
	quickStats := s.getQuickStats(personId, periodMonth, periodYear)

	// Get ranking
	ranking := s.getRanking(personId, periodMonth, periodYear)

	dashboard := dto.TLDashboardResponse{
		PersonInfo:        personInfo,
		CurrentEvaluation: currentEval,
		RecentActivities:  recentActivities,
		PerformanceTrend:  performanceTrend,
		QuickStats:        quickStats,
		Ranking:           ranking,
	}

	return dashboard, nil
}

// ========================================
// HELPER FUNCTIONS
// ========================================

func (s *TLDashboardService) getPersonInfo(personId string) (dto.PersonInfo, error) {
	var person domainperson.Person
	if err := s.DB.Where("id = ?", personId).First(&person).Error; err != nil {
		return dto.PersonInfo{}, err
	}

	return dto.PersonInfo{
		PersonId:   person.Id,
		Name:       person.Name,
		HondaId:    person.HondaId,
		DealerCode: person.DealerCode,
		Role:       person.Role,
	}, nil
}

func (s *TLDashboardService) buildEvaluationSummary(eval domainevaluation.Evaluation) (dto.EvaluationSummary, error) {
	// Get details
	details, err := s.EvalRepo.GetDetailsByEvaluationId(eval.Id)
	if err != nil {
		return dto.EvaluationSummary{}, err
	}

	// Calculate pillar scores
	pillarScores := s.calculatePillarScores(details)

	// Get top and weak KPIs
	topKPIs, weakKPIs := s.getTopAndWeakKPIs(details, 3)

	summary := dto.EvaluationSummary{
		EvaluationId: eval.Id,
		PeriodMonth:  eval.Period.PeriodMonth,
		PeriodYear:   eval.Period.PeriodYear,
		TotalScore:   eval.TotalScore,
		PillarScores: pillarScores,
		TopKPIs:      topKPIs,
		WeakKPIs:     weakKPIs,
		LastUpdated:  eval.CreatedAt,
	}

	return summary, nil
}

func (s *TLDashboardService) calculatePillarScores(details []domainevaluation.EvaluationDetail) []dto.PillarScoreSummary {
	// Group scores by pillar
	pillarScores := make(map[string]*dto.PillarScoreSummary)

	// Get all KPI items with pillar info
	type KPIPillar struct {
		KPIId      string
		PillarId   string
		PillarName string
		Weight     float64
	}

	var kpiPillars []KPIPillar
	s.DB.Table("kpi_items k").
		Select("k.id as kpi_id, k.pillar_id, p.name as pillar_name, p.weight").
		Joins("INNER JOIN pillars p ON k.pillar_id = p.id").
		Scan(&kpiPillars)

	kpiMap := make(map[string]KPIPillar)
	for _, kp := range kpiPillars {
		kpiMap[kp.KPIId] = kp
	}

	// Aggregate scores by pillar
	for _, detail := range details {
		kp, found := kpiMap[detail.KpiItemId]
		if !found {
			continue
		}

		if _, exists := pillarScores[kp.PillarName]; !exists {
			pillarScores[kp.PillarName] = &dto.PillarScoreSummary{
				Name:     kp.PillarName,
				Score:    0,
				MaxScore: kp.Weight,
			}
		}

		pillarScores[kp.PillarName].Score += detail.Score
	}

	// Calculate percentages
	var result []dto.PillarScoreSummary
	for _, ps := range pillarScores {
		if ps.MaxScore > 0 {
			ps.Percentage = (ps.Score / ps.MaxScore) * 100
		}
		result = append(result, *ps)
	}

	// Sort by score descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result
}

func (s *TLDashboardService) getTopAndWeakKPIs(details []domainevaluation.EvaluationDetail, limit int) ([]dto.KPISummary, []dto.KPISummary) {
	// Get KPI info
	type KPIInfo struct {
		Id     string
		Name   string
		Weight float64
	}

	var kpiInfos []KPIInfo
	s.DB.Table("kpi_items").Select("id, name, weight").Scan(&kpiInfos)
	kpiMap := make(map[string]KPIInfo)
	for _, ki := range kpiInfos {
		kpiMap[ki.Id] = ki
	}

	// Build KPI summaries
	var summaries []dto.KPISummary
	for _, detail := range details {
		ki, found := kpiMap[detail.KpiItemId]
		if !found {
			continue
		}

		percentage := float64(0)
		if ki.Weight > 0 {
			percentage = (detail.Score / ki.Weight) * 100
		}

		summary := dto.KPISummary{
			Name:        ki.Name,
			Score:       detail.Score,
			MaxScore:    ki.Weight,
			ActualValue: detail.ActualValue,
			Percentage:  percentage,
		}
		summaries = append(summaries, summary)
	}

	// Sort by percentage
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Percentage > summaries[j].Percentage
	})

	// Get top KPIs
	topLimit := limit
	if len(summaries) < topLimit {
		topLimit = len(summaries)
	}
	topKPIs := summaries[:topLimit]

	// Get weak KPIs (reverse order)
	weakLimit := limit
	if len(summaries) < weakLimit {
		weakLimit = len(summaries)
	}
	weakKPIs := make([]dto.KPISummary, weakLimit)
	for i := 0; i < weakLimit; i++ {
		weakKPIs[i] = summaries[len(summaries)-1-i]
	}

	return topKPIs, weakKPIs
}

func (s *TLDashboardService) getRecentActivities(personId string, limit int) []dto.RecentActivityItem {
	var activities []dto.RecentActivityItem

	// Get recent daily activities
	type Activity struct {
		Date         time.Time
		ActivityType string
		Notes        *string
	}

	var tlActivities []Activity
	s.DB.Table("tl_daily_activities").
		Select("date, activity_type, notes").
		Where("person_id = ? AND deleted_at IS NULL", personId).
		Order("date DESC").
		Limit(limit).
		Scan(&tlActivities)

	for _, a := range tlActivities {
		desc := fmt.Sprintf("%s activity", a.ActivityType)
		details := ""
		if a.Notes != nil {
			details = *a.Notes
		}

		activities = append(activities, dto.RecentActivityItem{
			Type:        "activity",
			Date:        a.Date,
			Description: desc,
			Details:     details,
		})
	}

	// Get recent sessions
	type Session struct {
		Date        time.Time
		SessionType string
		Notes       *string
	}

	var sessions []Session
	s.DB.Table("tl_sessions").
		Select("date, session_type, notes").
		Where("person_id = ? AND deleted_at IS NULL", personId).
		Order("date DESC").
		Limit(limit).
		Scan(&sessions)

	for _, sess := range sessions {
		desc := fmt.Sprintf("%s session", sess.SessionType)
		details := ""
		if sess.Notes != nil {
			details = *sess.Notes
		}

		activities = append(activities, dto.RecentActivityItem{
			Type:        sess.SessionType,
			Date:        sess.Date,
			Description: desc,
			Details:     details,
		})
	}

	// Sort by date descending
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Date.After(activities[j].Date)
	})

	// Limit results
	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities
}

func (s *TLDashboardService) getPerformanceTrend(personId string, months int) []dto.PerformanceTrendItem {
	var trends []dto.PerformanceTrendItem

	// Get evaluations for last N months
	type EvalTrend struct {
		PeriodMonth int
		PeriodYear  int
		TotalScore  float64
	}

	var evalTrends []EvalTrend
	s.DB.Table("evaluations e").
		Select("ep.period_month, ep.period_year, e.total_score").
		Joins("INNER JOIN evaluation_periods ep ON e.evaluation_period_id = ep.id").
		Where("e.person_id = ?", personId).
		Order("ep.period_year DESC, ep.period_month DESC").
		Limit(months).
		Scan(&evalTrends)

	// Reverse to chronological order
	for i := len(evalTrends) - 1; i >= 0; i-- {
		et := evalTrends[i]
		label := fmt.Sprintf("%s %d", time.Month(et.PeriodMonth).String()[:3], et.PeriodYear)

		trends = append(trends, dto.PerformanceTrendItem{
			PeriodMonth: et.PeriodMonth,
			PeriodYear:  et.PeriodYear,
			PeriodLabel: label,
			TotalScore:  et.TotalScore,
		})
	}

	return trends
}

func (s *TLDashboardService) getQuickStats(personId string, periodMonth int, periodYear int) dto.TLQuickStats {
	// Calculate period boundaries
	startDate := time.Date(periodYear, time.Month(periodMonth), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	// Count activities
	var activityCount int64
	s.DB.Table("tl_daily_activities").
		Where("person_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", personId, startDate, endDate).
		Count(&activityCount)

	// Count coaching sessions
	var coachingCount int64
	s.DB.Table("tl_sessions").
		Where("person_id = ? AND session_type = 'coaching' AND date >= ? AND date <= ? AND deleted_at IS NULL", personId, startDate, endDate).
		Count(&coachingCount)

	// Count briefing sessions
	var briefingCount int64
	s.DB.Table("tl_sessions").
		Where("person_id = ? AND session_type = 'briefing' AND date >= ? AND date <= ? AND deleted_at IS NULL", personId, startDate, endDate).
		Count(&briefingCount)

	// Count training participations
	var trainingCount int64
	s.DB.Table("tl_training_participations").
		Where("person_id = ? AND start_date >= ? AND start_date <= ? AND deleted_at IS NULL", personId, startDate, endDate).
		Count(&trainingCount)

	// Get team size (TODO: implement proper supervisor relationship)
	var teamSize int64
	s.DB.Table("persons").
		Where("role = 'salesman' AND active = true").
		Count(&teamSize)

	// Calculate attendance rate
	var result struct {
		AvgRate float64
	}
	s.DB.Table("tl_attendance_records").
		Select("COALESCE(AVG(CASE WHEN status = 'hadir' THEN 100.0 ELSE 0.0 END), 0) as avg_rate").
		Where("tl_person_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", personId, startDate, endDate).
		Scan(&result)

	return dto.TLQuickStats{
		ActivitiesCount: int(activityCount),
		CoachingCount:   int(coachingCount),
		BriefingCount:   int(briefingCount),
		TrainingCount:   int(trainingCount),
		TeamSize:        int(teamSize),
		AttendanceRate:  result.AvgRate,
	}
}

func (s *TLDashboardService) getRanking(personId string, periodMonth int, periodYear int) *dto.RankingInfo {
	period, err := s.EvalRepo.GetPeriodByMonthYear(periodMonth, periodYear)
	if err != nil {
		return nil
	}

	// Get all evaluations for this period
	evaluations, err := s.EvalRepo.GetLeaderboard(period.Id, 0) // 0 = no limit
	if err != nil {
		return nil
	}

	if len(evaluations) == 0 {
		return nil
	}

	// Find rank
	rank := 0
	for i, eval := range evaluations {
		if eval.PersonId == personId {
			rank = i + 1
			break
		}
	}

	if rank == 0 {
		return nil
	}

	// Calculate percentile
	percentile := int((float64(rank) / float64(len(evaluations))) * 100)

	return &dto.RankingInfo{
		Rank:       rank,
		TotalTLs:   len(evaluations),
		Percentile: 100 - percentile, // Top X%
	}
}
