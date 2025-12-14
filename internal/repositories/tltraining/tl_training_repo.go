package repositorytltraining

import (
	"fmt"
	domaintltraining "teamleader-management/internal/domain/tltraining"
	interfacetltraining "teamleader-management/internal/interfaces/tltraining"
	"teamleader-management/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewTLTrainingRepo(db *gorm.DB) interfacetltraining.RepoTLTrainingInterface {
	return &repo{DB: db}
}

func (r *repo) StoreMultiple(records []domaintltraining.TLTrainingParticipation) error {
	return r.DB.Create(&records).Error
}

func (r *repo) GetByTrainingBatch(trainingBatch string) ([]domaintltraining.TLTrainingParticipation, error) {
	var ret []domaintltraining.TLTrainingParticipation
	if err := r.DB.Where("training_batch = ?", trainingBatch).Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *repo) GetAll(params filter.BaseParams) ([]domaintltraining.TLTrainingParticipation, int64, error) {
	var (
		ret       []domaintltraining.TLTrainingParticipation
		totalData int64
	)

	query := r.DB.Model(&domaintltraining.TLTrainingParticipation{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(training_name) LIKE LOWER(?) OR LOWER(salesman_name) LIKE LOWER(?)", searchPattern, searchPattern)
	}

	for key, value := range params.Filters {
		if value == nil {
			continue
		}

		switch key {
		case "tl_person_id":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("tl_person_id = ?", s)
			}
		case "salesman_id":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("salesman_id = ?", s)
			}
		case "status":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("status = ?", s)
			}
		case "date_from":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("date >= ?", s)
			}
		case "date_to":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("date <= ?", s)
			}
		}
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	if params.OrderBy != "" && params.OrderDirection != "" {
		validColumns := map[string]bool{
			"id":            true,
			"tl_person_id":  true,
			"training_name": true,
			"date":          true,
			"created_at":    true,
		}

		if _, ok := validColumns[params.OrderBy]; !ok {
			return nil, 0, fmt.Errorf("invalid orderBy column: %s", params.OrderBy)
		}

		query = query.Order(fmt.Sprintf("%s %s", params.OrderBy, params.OrderDirection))
	}

	if err := query.Offset(params.Offset).Limit(params.Limit).Find(&ret).Error; err != nil {
		return nil, 0, err
	}

	return ret, totalData, nil
}

var _ interfacetltraining.RepoTLTrainingInterface = (*repo)(nil)
