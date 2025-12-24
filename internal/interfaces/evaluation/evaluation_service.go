package interfaceevaluation

import (
	"teamleader-management/internal/dto"
	"teamleader-management/pkg/filter"
)

type ServiceEvaluationInterface interface {
	// Calculate evaluation for a period (for one person or all persons)
	CalculateEvaluation(periodMonth int, periodYear int, personId string) ([]dto.EvaluationResponse, error)

	// Get evaluation by ID with full breakdown
	GetByID(id string) (dto.EvaluationResponse, error)

	// Get evaluation for a person in a specific period
	GetByPersonAndPeriod(personId string, periodMonth int, periodYear int) (dto.EvaluationResponse, error)

	// List evaluations with filters
	GetAll(params filter.BaseParams) ([]dto.EvaluationResponse, int64, error)

	// Get leaderboard for a period
	GetLeaderboard(periodMonth int, periodYear int, limit int) (dto.LeaderboardResponse, error)

	// Recalculate (delete old and calculate new)
	RecalculateEvaluation(periodMonth int, periodYear int, personId string) ([]dto.EvaluationResponse, error)
}
