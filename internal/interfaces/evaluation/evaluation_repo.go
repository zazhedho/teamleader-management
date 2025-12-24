package interfaceevaluation

import (
	domainevaluation "teamleader-management/internal/domain/evaluation"
	"teamleader-management/pkg/filter"
)

type RepoEvaluationInterface interface {
	// Evaluation CRUD
	Store(evaluation domainevaluation.Evaluation) error
	GetByID(id string) (domainevaluation.Evaluation, error)
	GetByPersonAndPeriod(personId string, periodId string) (domainevaluation.Evaluation, error)
	GetAll(params filter.BaseParams) ([]domainevaluation.Evaluation, int64, error)
	Update(evaluation domainevaluation.Evaluation) error
	Delete(id string) error

	// Evaluation Details
	StoreDetails(details []domainevaluation.EvaluationDetail) error
	GetDetailsByEvaluationId(evaluationId string) ([]domainevaluation.EvaluationDetail, error)
	DeleteDetailsByEvaluationId(evaluationId string) error

	// Evaluation Period
	StorePeriod(period domainevaluation.EvaluationPeriod) error
	GetPeriodByMonthYear(month int, year int) (domainevaluation.EvaluationPeriod, error)
	GetOrCreatePeriod(month int, year int) (domainevaluation.EvaluationPeriod, error)

	// Leaderboard
	GetLeaderboard(periodId string, limit int) ([]domainevaluation.Evaluation, error)
}
