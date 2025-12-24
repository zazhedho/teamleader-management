package repositoryevaluation

import (
	domainevaluation "teamleader-management/internal/domain/evaluation"
	interfaceevaluation "teamleader-management/internal/interfaces/evaluation"
	"teamleader-management/pkg/filter"
	"teamleader-management/utils"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewEvaluationRepo(db *gorm.DB) interfaceevaluation.RepoEvaluationInterface {
	return &repo{DB: db}
}

// ========================================
// Evaluation CRUD
// ========================================

func (r *repo) Store(evaluation domainevaluation.Evaluation) error {
	return r.DB.Create(&evaluation).Error
}

func (r *repo) GetByID(id string) (domainevaluation.Evaluation, error) {
	var evaluation domainevaluation.Evaluation
	err := r.DB.Preload("Period").
		Preload("Details").
		Where("id = ?", id).
		First(&evaluation).Error
	if err != nil {
		return domainevaluation.Evaluation{}, err
	}
	return evaluation, nil
}

func (r *repo) GetByPersonAndPeriod(personId string, periodId string) (domainevaluation.Evaluation, error) {
	var evaluation domainevaluation.Evaluation
	err := r.DB.Preload("Period").
		Preload("Details").
		Where("person_id = ? AND evaluation_period_id = ?", personId, periodId).
		First(&evaluation).Error
	if err != nil {
		return domainevaluation.Evaluation{}, err
	}
	return evaluation, nil
}

func (r *repo) GetAll(params filter.BaseParams) ([]domainevaluation.Evaluation, int64, error) {
	var evaluations []domainevaluation.Evaluation
	var total int64

	query := r.DB.Model(&domainevaluation.Evaluation{})

	// Apply filters
	if periodId, ok := params.Filters["evaluation_period_id"].(string); ok && periodId != "" {
		query = query.Where("evaluation_period_id = ?", periodId)
	}

	if personId, ok := params.Filters["person_id"].(string); ok && personId != "" {
		query = query.Where("person_id = ?", personId)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if params.OrderBy != "" {
		query = query.Order(params.OrderBy + " " + params.OrderDirection)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute query with preloads
	err := query.Preload("Period").Preload("Details").Find(&evaluations).Error
	if err != nil {
		return nil, 0, err
	}

	return evaluations, total, nil
}

func (r *repo) Update(evaluation domainevaluation.Evaluation) error {
	return r.DB.Save(&evaluation).Error
}

func (r *repo) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&domainevaluation.Evaluation{}).Error
}

// ========================================
// Evaluation Details
// ========================================

func (r *repo) StoreDetails(details []domainevaluation.EvaluationDetail) error {
	if len(details) == 0 {
		return nil
	}
	return r.DB.Create(&details).Error
}

func (r *repo) GetDetailsByEvaluationId(evaluationId string) ([]domainevaluation.EvaluationDetail, error) {
	var details []domainevaluation.EvaluationDetail
	err := r.DB.Where("evaluation_id = ?", evaluationId).
		Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (r *repo) DeleteDetailsByEvaluationId(evaluationId string) error {
	return r.DB.Where("evaluation_id = ?", evaluationId).
		Delete(&domainevaluation.EvaluationDetail{}).Error
}

// ========================================
// Evaluation Period
// ========================================

func (r *repo) StorePeriod(period domainevaluation.EvaluationPeriod) error {
	return r.DB.Create(&period).Error
}

func (r *repo) GetPeriodByMonthYear(month int, year int) (domainevaluation.EvaluationPeriod, error) {
	var period domainevaluation.EvaluationPeriod
	err := r.DB.Where("period_month = ? AND period_year = ?", month, year).
		First(&period).Error
	if err != nil {
		return domainevaluation.EvaluationPeriod{}, err
	}
	return period, nil
}

func (r *repo) GetOrCreatePeriod(month int, year int) (domainevaluation.EvaluationPeriod, error) {
	// Try to get existing
	period, err := r.GetPeriodByMonthYear(month, year)
	if err == nil {
		return period, nil
	}

	// If not found, create new
	if err == gorm.ErrRecordNotFound {
		newPeriod := domainevaluation.EvaluationPeriod{
			Id:          utils.CreateUUID(),
			PeriodMonth: month,
			PeriodYear:  year,
		}
		if err := r.StorePeriod(newPeriod); err != nil {
			return domainevaluation.EvaluationPeriod{}, err
		}
		return newPeriod, nil
	}

	return domainevaluation.EvaluationPeriod{}, err
}

// ========================================
// Leaderboard
// ========================================

func (r *repo) GetLeaderboard(periodId string, limit int) ([]domainevaluation.Evaluation, error) {
	var evaluations []domainevaluation.Evaluation
	query := r.DB.Where("evaluation_period_id = ?", periodId).
		Order("total_score DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Preload("Period").Find(&evaluations).Error
	if err != nil {
		return nil, err
	}

	return evaluations, nil
}

var _ interfaceevaluation.RepoEvaluationInterface = (*repo)(nil)
