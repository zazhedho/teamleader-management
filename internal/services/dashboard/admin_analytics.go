package servicedashboard

import (
	"fmt"
	"math"
	"sort"
	"time"

	domainevaluation "teamleader-management/internal/domain/evaluation"
	domainperson "teamleader-management/internal/domain/person"
	"teamleader-management/internal/dto"
	interfaceevaluation "teamleader-management/internal/interfaces/evaluation"

	"gorm.io/gorm"
)

type AdminAnalyticsService struct {
	DB       *gorm.DB
	EvalRepo interfaceevaluation.RepoEvaluationInterface
}

func NewAdminAnalyticsService(db *gorm.DB, evalRepo interfaceevaluation.RepoEvaluationInterface) *AdminAnalyticsService {
	return &AdminAnalyticsService{
		DB:       db,
		EvalRepo: evalRepo,
	}
}

// GetAnalytics retrieves comprehensive analytics for admin
func (s *AdminAnalyticsService) GetAnalytics(periodMonth int, periodYear int, topN int, trendMonths int) (dto.AdminAnalyticsResponse, error) {
	period, err := s.EvalRepo.GetPeriodByMonthYear(periodMonth, periodYear)
	if err != nil {
		return dto.AdminAnalyticsResponse{}, fmt.Errorf("period not found: %w", err)
	}

	// Get all evaluations for this period
	evaluations, err := s.EvalRepo.GetLeaderboard(period.Id, 0) // 0 = all
	if err != nil {
		return dto.AdminAnalyticsResponse{}, err
	}

	if len(evaluations) == 0 {
		return dto.AdminAnalyticsResponse{}, fmt.Errorf("no evaluations found for this period")
	}

	// Overall statistics
	overallStats := s.calculateOverallStats(evaluations)

	// Top and bottom performers
	topPerformers := s.getTopPerformers(evaluations, topN)
	bottomPerformers := s.getBottomPerformers(evaluations, topN)

	// Pillar analysis
	pillarAnalysis := s.getPillarAnalysis(period.Id)

	// Trend comparison
	trendComparison := s.getTrendComparison(periodMonth, periodYear, trendMonths)

	// Score distributions
	distributions := s.getScoreDistribution(evaluations)

	analytics := dto.AdminAnalyticsResponse{
		Period:           fmt.Sprintf("%d-%02d", periodYear, periodMonth),
		OverallStats:     overallStats,
		TopPerformers:    topPerformers,
		BottomPerformers: bottomPerformers,
		PillarAnalysis:   pillarAnalysis,
		TrendComparison:  trendComparison,
		Distributions:    distributions,
	}

	return analytics, nil
}

// CompareTeam compares multiple TLs side by side
func (s *AdminAnalyticsService) CompareTeam(personIds []string, periodMonth int, periodYear int) (dto.TeamComparisonResponse, error) {
	period, err := s.EvalRepo.GetPeriodByMonthYear(periodMonth, periodYear)
	if err != nil {
		return dto.TeamComparisonResponse{}, fmt.Errorf("period not found: %w", err)
	}

	var comparisons []dto.TLComparisonItem
	var allPillarScores []dto.PillarComparisonChart

	// Get person names
	var persons []domainperson.Person
	s.DB.Where("id IN ?", personIds).Find(&persons)
	personMap := make(map[string]string)
	for _, p := range persons {
		personMap[p.Id] = p.Name
	}

	// Get evaluations for each person
	for rank, personId := range personIds {
		eval, err := s.EvalRepo.GetByPersonAndPeriod(personId, period.Id)
		if err != nil {
			continue
		}

		// Get details
		details, _ := s.EvalRepo.GetDetailsByEvaluationId(eval.Id)

		// Calculate pillar scores
		pillarScores := s.calculatePillarScoresForComparison(details)

		comparison := dto.TLComparisonItem{
			PersonId:     personId,
			PersonName:   personMap[personId],
			TotalScore:   eval.TotalScore,
			Rank:         rank + 1,
			PillarScores: pillarScores,
		}
		comparisons = append(comparisons, comparison)

		// Collect for chart
		for _, ps := range pillarScores {
			s.addToPillarChart(&allPillarScores, ps.Name, personMap[personId], ps.Score)
		}
	}

	// Sort by total score descending
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].TotalScore > comparisons[j].TotalScore
	})

	// Update ranks
	for i := range comparisons {
		comparisons[i].Rank = i + 1
	}

	response := dto.TeamComparisonResponse{
		Period:      fmt.Sprintf("%d-%02d", periodYear, periodMonth),
		Comparisons: comparisons,
		PillarChart: allPillarScores,
	}

	return response, nil
}

// ========================================
// HELPER FUNCTIONS
// ========================================

func (s *AdminAnalyticsService) calculateOverallStats(evaluations []domainevaluation.Evaluation) dto.OverallStatistics {
	scores := make([]float64, len(evaluations))
	sum := 0.0
	min := 100.0
	max := 0.0

	for i, eval := range evaluations {
		score := eval.TotalScore
		scores[i] = score
		sum += score

		if score < min {
			min = score
		}
		if score > max {
			max = score
		}
	}

	count := float64(len(scores))
	average := sum / count

	// Calculate median
	sortedScores := make([]float64, len(scores))
	copy(sortedScores, scores)
	sort.Float64s(sortedScores)

	median := 0.0
	if len(sortedScores)%2 == 0 {
		median = (sortedScores[len(sortedScores)/2-1] + sortedScores[len(sortedScores)/2]) / 2
	} else {
		median = sortedScores[len(sortedScores)/2]
	}

	// Calculate standard deviation
	variance := 0.0
	for _, score := range scores {
		variance += math.Pow(score-average, 2)
	}
	stdDev := math.Sqrt(variance / count)

	return dto.OverallStatistics{
		TotalTLs:          len(evaluations),
		AverageScore:      average,
		MedianScore:       median,
		HighestScore:      max,
		LowestScore:       min,
		StandardDeviation: stdDev,
	}
}

func (s *AdminAnalyticsService) getTopPerformers(evaluations []domainevaluation.Evaluation, topN int) []dto.LeaderboardEntry {
	// Already sorted by score descending from GetLeaderboard
	limit := topN
	if len(evaluations) < limit {
		limit = len(evaluations)
	}

	return s.buildLeaderboardEntries(evaluations[:limit])
}

func (s *AdminAnalyticsService) getBottomPerformers(evaluations []domainevaluation.Evaluation, bottomN int) []dto.LeaderboardEntry {
	limit := bottomN
	if len(evaluations) < limit {
		limit = len(evaluations)
	}

	// Get last N items (reverse order)
	bottomEvals := make([]domainevaluation.Evaluation, limit)
	for i := 0; i < limit; i++ {
		bottomEvals[i] = evaluations[len(evaluations)-1-i]
	}

	return s.buildLeaderboardEntries(bottomEvals)
}

func (s *AdminAnalyticsService) buildLeaderboardEntries(evaluations []domainevaluation.Evaluation) []dto.LeaderboardEntry {
	var entries []dto.LeaderboardEntry

	for rank, eval := range evaluations {
		personId := eval.PersonId

		// Get person info
		var person domainperson.Person
		s.DB.Where("id = ?", personId).First(&person)

		dealerCode := ""
		if person.DealerCode != nil {
			dealerCode = *person.DealerCode
		}

		entry := dto.LeaderboardEntry{
			Rank:        rank + 1,
			PersonId:    personId,
			PersonName:  person.Name,
			DealerCode:  dealerCode,
			TotalScore:  eval.TotalScore,
			PeriodMonth: eval.Period.PeriodMonth,
			PeriodYear:  eval.Period.PeriodYear,
		}
		entries = append(entries, entry)
	}

	return entries
}

func (s *AdminAnalyticsService) getPillarAnalysis(periodId string) []dto.PillarAnalysis {
	type PillarStats struct {
		PillarName   string
		AverageScore float64
		MaxPossible  float64
		TopScorer    string
		TopScore     float64
	}

	var pillarAnalyses []dto.PillarAnalysis

	// Get all pillars
	type Pillar struct {
		Id     string
		Name   string
		Weight float64
	}

	var pillars []Pillar
	s.DB.Table("pillars").Select("id, name, weight").Scan(&pillars)

	for _, pillar := range pillars {
		// Get average score for this pillar
		var result struct {
			AvgScore float64
			TopScore float64
		}

		s.DB.Table("evaluation_details ed").
			Select("COALESCE(AVG(ed.score), 0) as avg_score, COALESCE(MAX(ed.score), 0) as top_score").
			Joins("INNER JOIN kpi_items k ON ed.kpi_item_id = k.id").
			Joins("INNER JOIN evaluations e ON ed.evaluation_id = e.id").
			Where("k.pillar_id = ? AND e.evaluation_period_id = ?", pillar.Id, periodId).
			Scan(&result)

		// Find top scorer for this pillar
		type TopScorer struct {
			PersonId   string
			PersonName string
			Score      float64
		}

		var topScorer TopScorer
		s.DB.Table("evaluation_details ed").
			Select("e.person_id, p.name as person_name, SUM(ed.score) as score").
			Joins("INNER JOIN kpi_items k ON ed.kpi_item_id = k.id").
			Joins("INNER JOIN evaluations e ON ed.evaluation_id = e.id").
			Joins("INNER JOIN persons p ON e.person_id = p.id").
			Where("k.pillar_id = ? AND e.evaluation_period_id = ?", pillar.Id, periodId).
			Group("e.person_id, p.name").
			Order("score DESC").
			Limit(1).
			Scan(&topScorer)

		achievementPct := (result.AvgScore / pillar.Weight) * 100

		analysis := dto.PillarAnalysis{
			PillarName:     pillar.Name,
			AverageScore:   result.AvgScore,
			MaxPossible:    pillar.Weight,
			AchievementPct: achievementPct,
			TopScorer:      topScorer.PersonName,
			TopScore:       topScorer.Score,
		}
		pillarAnalyses = append(pillarAnalyses, analysis)
	}

	return pillarAnalyses
}

func (s *AdminAnalyticsService) getTrendComparison(currentMonth int, currentYear int, months int) []dto.PeriodComparison {
	var trends []dto.PeriodComparison

	// Get last N months of data
	for i := months - 1; i >= 0; i-- {
		targetDate := time.Date(currentYear, time.Month(currentMonth), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -i, 0)
		month := int(targetDate.Month())
		year := targetDate.Year()

		// Get period
		period, err := s.EvalRepo.GetPeriodByMonthYear(month, year)
		if err != nil {
			continue
		}

		// Get evaluations for this period
		evaluations, err := s.EvalRepo.GetLeaderboard(period.Id, 0)
		if err != nil {
			continue
		}

		if len(evaluations) == 0 {
			continue
		}

		// Calculate average
		sum := 0.0
		for _, eval := range evaluations {
			sum += eval.TotalScore
		}
		avgScore := sum / float64(len(evaluations))

		label := fmt.Sprintf("%s %d", time.Month(month).String()[:3], year)

		// Calculate change vs previous
		var change *float64
		if len(trends) > 0 {
			prevAvg := trends[len(trends)-1].AverageScore
			diff := avgScore - prevAvg
			change = &diff
		}

		trend := dto.PeriodComparison{
			PeriodMonth:      month,
			PeriodYear:       year,
			PeriodLabel:      label,
			AverageScore:     avgScore,
			ParticipatingTLs: len(evaluations),
			Change:           change,
		}
		trends = append(trends, trend)
	}

	return trends
}

func (s *AdminAnalyticsService) getScoreDistribution(evaluations []domainevaluation.Evaluation) dto.ScoreDistribution {
	dist := dto.ScoreDistribution{}

	for _, eval := range evaluations {
		score := eval.TotalScore

		switch {
		case score >= 0 && score < 20:
			dist.Range0_20++
		case score >= 20 && score < 40:
			dist.Range20_40++
		case score >= 40 && score < 60:
			dist.Range40_60++
		case score >= 60 && score < 80:
			dist.Range60_80++
		case score >= 80 && score <= 100:
			dist.Range80_100++
		}
	}

	return dist
}

func (s *AdminAnalyticsService) calculatePillarScoresForComparison(details []domainevaluation.EvaluationDetail) []dto.PillarScoreSummary {
	// This is similar to TL dashboard service but simplified
	pillarScores := make(map[string]*dto.PillarScoreSummary)

	// Get KPI-Pillar mapping
	type KPIPillar struct {
		KPIId      string
		PillarName string
		Weight     float64
	}

	var kpiPillars []KPIPillar
	s.DB.Table("kpi_items k").
		Select("k.id as kpi_id, p.name as pillar_name, p.weight").
		Joins("INNER JOIN pillars p ON k.pillar_id = p.id").
		Scan(&kpiPillars)

	kpiMap := make(map[string]KPIPillar)
	for _, kp := range kpiPillars {
		kpiMap[kp.KPIId] = kp
	}

	// Aggregate
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

	// Convert to slice and calculate percentages
	var result []dto.PillarScoreSummary
	for _, ps := range pillarScores {
		if ps.MaxScore > 0 {
			ps.Percentage = (ps.Score / ps.MaxScore) * 100
		}
		result = append(result, *ps)
	}

	return result
}

func (s *AdminAnalyticsService) addToPillarChart(charts *[]dto.PillarComparisonChart, pillarName string, personName string, score float64) {
	// Find or create pillar chart
	var chart *dto.PillarComparisonChart
	for i := range *charts {
		if (*charts)[i].PillarName == pillarName {
			chart = &(*charts)[i]
			break
		}
	}

	if chart == nil {
		*charts = append(*charts, dto.PillarComparisonChart{
			PillarName: pillarName,
			TLScores:   []dto.TLPillarScore{},
		})
		chart = &(*charts)[len(*charts)-1]
	}

	// Add score
	chart.TLScores = append(chart.TLScores, dto.TLPillarScore{
		PersonName: personName,
		Score:      score,
	})
}
