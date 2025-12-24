package serviceevaluation

import (
	"fmt"
	"time"

	domainevaluation "teamleader-management/internal/domain/evaluation"
	domainkpiitem "teamleader-management/internal/domain/kpiitem"
	domainperson "teamleader-management/internal/domain/person"
	domainpillar "teamleader-management/internal/domain/pillar"
	"teamleader-management/internal/dto"
	interfaceevaluation "teamleader-management/internal/interfaces/evaluation"
	"teamleader-management/pkg/filter"
	"teamleader-management/utils"

	"gorm.io/gorm"
)

type ServiceEvaluation struct {
	Repo       interfaceevaluation.RepoEvaluationInterface
	Calculator *EvaluationCalculator
	DB         *gorm.DB
}

func NewEvaluationService(repo interfaceevaluation.RepoEvaluationInterface, db *gorm.DB) *ServiceEvaluation {
	return &ServiceEvaluation{
		Repo:       repo,
		Calculator: NewEvaluationCalculator(db),
		DB:         db,
	}
}

// CalculateEvaluation calculates evaluation for a period
// If personId is empty, calculates for all TLs
func (s *ServiceEvaluation) CalculateEvaluation(periodMonth int, periodYear int, personId string) ([]dto.EvaluationResponse, error) {
	// 1. Get or create evaluation period
	period, err := s.Repo.GetOrCreatePeriod(periodMonth, periodYear)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create period: %w", err)
	}

	// 2. Get list of TLs to evaluate
	tlPersonIds, err := s.getTLPersonIds(personId)
	if err != nil {
		return nil, fmt.Errorf("failed to get TL person IDs: %w", err)
	}

	if len(tlPersonIds) == 0 {
		return nil, fmt.Errorf("no team leaders found for evaluation")
	}

	// 3. Calculate evaluation for each TL
	var results []dto.EvaluationResponse

	for _, tlPersonId := range tlPersonIds {
		// Calculate scores
		calculationResult, err := s.Calculator.Calculate(tlPersonId, periodMonth, periodYear)
		if err != nil {
			// Log error but continue with other TLs
			continue
		}

		// Check if evaluation already exists
		existing, err := s.Repo.GetByPersonAndPeriod(tlPersonId, period.Id)
		if err == nil {
			// Evaluation exists, delete old details and update
			_ = s.Repo.DeleteDetailsByEvaluationId(existing.Id)
			existing.TotalScore = calculationResult.TotalScore
			if err := s.Repo.Update(existing); err != nil {
				continue
			}

			// Store new details
			for i := range calculationResult.Details {
				calculationResult.Details[i].EvaluationId = existing.Id
			}
			if err := s.Repo.StoreDetails(calculationResult.Details); err != nil {
				continue
			}

			// Get full result with details
			response, err := s.buildEvaluationResponse(existing.Id)
			if err != nil {
				continue
			}
			results = append(results, response)

		} else {
			// Create new evaluation
			evaluation := domainevaluation.Evaluation{
				Id:                 utils.CreateUUID(),
				EvaluationPeriodId: period.Id,
				PersonId:           tlPersonId,
				TotalScore:         calculationResult.TotalScore,
				CreatedAt:          time.Now(),
			}

			if err := s.Repo.Store(evaluation); err != nil {
				continue
			}

			// Store details
			for i := range calculationResult.Details {
				calculationResult.Details[i].EvaluationId = evaluation.Id
			}
			if err := s.Repo.StoreDetails(calculationResult.Details); err != nil {
				continue
			}

			// Get full result with details
			response, err := s.buildEvaluationResponse(evaluation.Id)
			if err != nil {
				continue
			}
			results = append(results, response)
		}
	}

	return results, nil
}

// GetByID retrieves evaluation with full breakdown
func (s *ServiceEvaluation) GetByID(id string) (dto.EvaluationResponse, error) {
	return s.buildEvaluationResponse(id)
}

// GetByPersonAndPeriod retrieves evaluation for a person in a specific period
func (s *ServiceEvaluation) GetByPersonAndPeriod(personId string, periodMonth int, periodYear int) (dto.EvaluationResponse, error) {
	period, err := s.Repo.GetPeriodByMonthYear(periodMonth, periodYear)
	if err != nil {
		return dto.EvaluationResponse{}, fmt.Errorf("period not found: %w", err)
	}

	evaluation, err := s.Repo.GetByPersonAndPeriod(personId, period.Id)
	if err != nil {
		return dto.EvaluationResponse{}, fmt.Errorf("evaluation not found: %w", err)
	}

	return s.buildEvaluationResponse(evaluation.Id)
}

// GetAll lists evaluations with filters
func (s *ServiceEvaluation) GetAll(params filter.BaseParams) ([]dto.EvaluationResponse, int64, error) {
	evaluations, total, err := s.Repo.GetAll(params)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.EvaluationResponse
	for _, eval := range evaluations {
		response, err := s.buildEvaluationResponse(eval.Id)
		if err != nil {
			continue
		}
		responses = append(responses, response)
	}

	return responses, total, nil
}

// GetLeaderboard retrieves top TLs for a period
func (s *ServiceEvaluation) GetLeaderboard(periodMonth int, periodYear int, limit int) (dto.LeaderboardResponse, error) {
	period, err := s.Repo.GetPeriodByMonthYear(periodMonth, periodYear)
	if err != nil {
		return dto.LeaderboardResponse{}, fmt.Errorf("period not found: %w", err)
	}

	evaluations, err := s.Repo.GetLeaderboard(period.Id, limit)
	if err != nil {
		return dto.LeaderboardResponse{}, err
	}

	// Build leaderboard entries
	var entries []dto.LeaderboardEntry
	for rank, eval := range evaluations {
		// Get person details
		var person domainperson.Person
		if err := s.DB.Where("id = ?", eval.PersonId).First(&person).Error; err != nil {
			continue
		}

		dealerCode := ""
		if person.DealerCode != nil {
			dealerCode = *person.DealerCode
		}

		entry := dto.LeaderboardEntry{
			Rank:        rank + 1,
			PersonId:    eval.PersonId,
			PersonName:  person.Name,
			DealerCode:  dealerCode,
			TotalScore:  eval.TotalScore,
			PeriodMonth: period.PeriodMonth,
			PeriodYear:  period.PeriodYear,
		}
		entries = append(entries, entry)
	}

	response := dto.LeaderboardResponse{
		Period:  fmt.Sprintf("%d-%02d", periodYear, periodMonth),
		Entries: entries,
		Total:   len(entries),
	}

	return response, nil
}

// RecalculateEvaluation deletes existing evaluation and recalculates
func (s *ServiceEvaluation) RecalculateEvaluation(periodMonth int, periodYear int, personId string) ([]dto.EvaluationResponse, error) {
	// Get period
	period, err := s.Repo.GetPeriodByMonthYear(periodMonth, periodYear)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	// Delete existing evaluations if period exists
	if err == nil {
		tlPersonIds, _ := s.getTLPersonIds(personId)
		for _, tlPersonId := range tlPersonIds {
			existing, err := s.Repo.GetByPersonAndPeriod(tlPersonId, period.Id)
			if err == nil {
				_ = s.Repo.DeleteDetailsByEvaluationId(existing.Id)
				_ = s.Repo.Delete(existing.Id)
			}
		}
	}

	// Recalculate
	return s.CalculateEvaluation(periodMonth, periodYear, personId)
}

// ========================================
// HELPER FUNCTIONS
// ========================================

// buildEvaluationResponse builds full response with breakdown
func (s *ServiceEvaluation) buildEvaluationResponse(evaluationId string) (dto.EvaluationResponse, error) {
	// Get evaluation
	evaluation, err := s.Repo.GetByID(evaluationId)
	if err != nil {
		return dto.EvaluationResponse{}, err
	}

	// Get person
	var person domainperson.Person
	if err := s.DB.Where("id = ?", evaluation.PersonId).First(&person).Error; err != nil {
		return dto.EvaluationResponse{}, err
	}

	// Get details
	details, err := s.Repo.GetDetailsByEvaluationId(evaluationId)
	if err != nil {
		return dto.EvaluationResponse{}, err
	}

	// Get KPI items for breakdown
	var kpiItems []domainkpiitem.KPIItem
	if err := s.DB.Find(&kpiItems).Error; err != nil {
		return dto.EvaluationResponse{}, err
	}

	// Get pillars
	var pillars []domainpillar.Pillar
	if err := s.DB.Find(&pillars).Error; err != nil {
		return dto.EvaluationResponse{}, err
	}
	pillarMap := make(map[string]domainpillar.Pillar)
	for _, p := range pillars {
		pillarMap[p.Id] = p
	}

	// Build KPI breakdown
	kpiBreakdown := s.buildKPIBreakdown(details, kpiItems, pillarMap)

	// Build pillar breakdown
	pillarBreakdown := s.buildPillarBreakdown(kpiBreakdown)

	response := dto.EvaluationResponse{
		Id:              evaluation.Id,
		PersonId:        evaluation.PersonId,
		PersonName:      person.Name,
		PeriodMonth:     evaluation.Period.PeriodMonth,
		PeriodYear:      evaluation.Period.PeriodYear,
		TotalScore:      evaluation.TotalScore,
		PillarBreakdown: pillarBreakdown,
		KpiBreakdown:    kpiBreakdown,
		CreatedAt:       evaluation.CreatedAt,
	}

	return response, nil
}

// buildKPIBreakdown creates detailed KPI breakdown
func (s *ServiceEvaluation) buildKPIBreakdown(details []domainevaluation.EvaluationDetail, kpiItems []domainkpiitem.KPIItem, pillarMap map[string]domainpillar.Pillar) []dto.KpiScoreBreakdown {
	var breakdown []dto.KpiScoreBreakdown

	// Map KPI items by ID for quick lookup
	kpiMap := make(map[string]domainkpiitem.KPIItem)
	for _, kpi := range kpiItems {
		kpiMap[kpi.Id] = kpi
	}

	for _, detail := range details {
		kpi, found := kpiMap[detail.KpiItemId]
		if !found {
			continue
		}

		// Get pillar name from pillarMap
		pillar, pillarFound := pillarMap[kpi.PillarId]
		pillarName := ""
		if pillarFound {
			pillarName = pillar.Name
		}

		item := dto.KpiScoreBreakdown{
			KpiItemId:        detail.KpiItemId,
			KpiItemName:      kpi.Name,
			PillarName:       pillarName,
			Weight:           kpi.Weight,
			ActualValue:      detail.ActualValue,
			TargetValue:      kpi.TargetValue,
			AchievementRatio: detail.AchievementRatio,
			Score:            detail.Score,
			MaxScore:         kpi.Weight,
			Unit:             kpi.Unit,
			InputSource:      kpi.InputSource,
		}
		breakdown = append(breakdown, item)
	}

	return breakdown
}

// buildPillarBreakdown aggregates KPI scores by pillar
func (s *ServiceEvaluation) buildPillarBreakdown(kpiBreakdown []dto.KpiScoreBreakdown) []dto.PillarScoreBreakdown {
	// Group by pillar
	pillarScores := make(map[string]*dto.PillarScoreBreakdown)

	// Get all pillars first
	var pillars []domainpillar.Pillar
	s.DB.Find(&pillars)
	pillarMap := make(map[string]domainpillar.Pillar)
	for _, p := range pillars {
		pillarMap[p.Name] = p
	}

	for _, kpi := range kpiBreakdown {
		if _, exists := pillarScores[kpi.PillarName]; !exists {
			// Get pillar from map
			pillar, found := pillarMap[kpi.PillarName]
			if !found {
				continue
			}

			pillarScores[kpi.PillarName] = &dto.PillarScoreBreakdown{
				PillarId:       pillar.Id,
				PillarName:     kpi.PillarName,
				PillarWeight:   pillar.Weight,
				PillarScore:    0,
				PillarMaxScore: pillar.Weight,
			}
		}

		pillarScores[kpi.PillarName].PillarScore += kpi.Score
	}

	// Convert map to slice
	var breakdown []dto.PillarScoreBreakdown
	for _, pillar := range pillarScores {
		breakdown = append(breakdown, *pillar)
	}

	return breakdown
}

// getTLPersonIds returns list of TL person IDs to evaluate
func (s *ServiceEvaluation) getTLPersonIds(personId string) ([]string, error) {
	if personId != "" {
		// Verify person is a TL
		var person domainperson.Person
		if err := s.DB.Where("id = ? AND role = ?", personId, utils.RoleTL).First(&person).Error; err != nil {
			return nil, fmt.Errorf("person not found or not a team leader: %w", err)
		}
		return []string{personId}, nil
	}

	// Get all active TLs
	var persons []domainperson.Person
	if err := s.DB.Where("role = ? AND active = ?", utils.RoleTL, true).Find(&persons).Error; err != nil {
		return nil, err
	}

	var personIds []string
	for _, p := range persons {
		personIds = append(personIds, p.Id)
	}

	return personIds, nil
}

var _ interfaceevaluation.ServiceEvaluationInterface = (*ServiceEvaluation)(nil)
