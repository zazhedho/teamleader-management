package repositorykpiitem

import (
	"fmt"

	domainkpiitem "teamleader-management/internal/domain/kpiitem"
	interfacekpiitem "teamleader-management/internal/interfaces/kpiitem"
	"teamleader-management/pkg/filter"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repo struct {
	DB *gorm.DB
}

func NewKPIItemRepo(db *gorm.DB) interfacekpiitem.RepoKPIItemInterface {
	return &repo{DB: db}
}

func (r *repo) Store(m domainkpiitem.KPIItem) error {
	return r.DB.Create(&m).Error
}

func (r *repo) GetByID(id string) (domainkpiitem.KPIItem, error) {
	var ret domainkpiitem.KPIItem
	if err := r.DB.Where("id = ?", id).First(&ret).Error; err != nil {
		return domainkpiitem.KPIItem{}, err
	}
	return ret, nil
}

func (r *repo) GetByNameAndPillar(name, pillarId string) (domainkpiitem.KPIItem, error) {
	var ret domainkpiitem.KPIItem
	if err := r.DB.Where("LOWER(name) = LOWER(?) AND pillar_id = ?", name, pillarId).First(&ret).Error; err != nil {
		return domainkpiitem.KPIItem{}, err
	}
	return ret, nil
}

func (r *repo) GetAll(params filter.BaseParams) ([]domainkpiitem.KPIItem, int64, error) {
	var (
		ret       []domainkpiitem.KPIItem
		totalData int64
	)

	query := r.DB.Model(&domainkpiitem.KPIItem{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?)", searchPattern)
	}

	for key, value := range params.Filters {
		switch key {
		case "pillar_id":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("pillar_id = ?", s)
			}
		case "input_source":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("input_source = ?", s)
			}
		case "applies_to_tl":
			query = query.Where("applies_to_tl = ?", value)
		case "applies_to_salesman":
			query = query.Where("applies_to_salesman = ?", value)
		}
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	if params.OrderBy != "" && params.OrderDirection != "" {
		validColumns := map[string]bool{
			"name":         true,
			"weight":       true,
			"created_at":   true,
			"updated_at":   true,
			"input_source": true,
		}
		if !validColumns[params.OrderBy] {
			return nil, 0, fmt.Errorf("invalid orderBy column: %s", params.OrderBy)
		}
		query = query.Order(fmt.Sprintf("%s %s", params.OrderBy, params.OrderDirection))
	}

	if err := query.Offset(params.Offset).Limit(params.Limit).Find(&ret).Error; err != nil {
		return nil, 0, err
	}

	return ret, totalData, nil
}

func (r *repo) Update(m domainkpiitem.KPIItem) error {
	return r.DB.Save(&m).Error
}

func (r *repo) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&domainkpiitem.KPIItem{}).Error
}

type personTargetRepo struct {
	DB *gorm.DB
}

func NewPersonKPITargetRepo(db *gorm.DB) interfacekpiitem.RepoPersonKPITargetInterface {
	return &personTargetRepo{DB: db}
}

func (r *personTargetRepo) Upsert(target domainkpiitem.PersonKPITarget) error {
	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "person_id"}, {Name: "kpi_item_id"}, {Name: "period_month"}, {Name: "period_year"}},
		DoUpdates: clause.AssignmentColumns([]string{"target_value", "last_modified", "updated_by"}),
	}).Create(&target).Error
}

func (r *personTargetRepo) Get(personId, kpiItemId string, periodMonth, periodYear int) (domainkpiitem.PersonKPITarget, error) {
	var ret domainkpiitem.PersonKPITarget
	if err := r.DB.Where("person_id = ? AND kpi_item_id = ? AND period_month = ? AND period_year = ?", personId, kpiItemId, periodMonth, periodYear).First(&ret).Error; err != nil {
		return domainkpiitem.PersonKPITarget{}, err
	}
	return ret, nil
}

func (r *personTargetRepo) Delete(personId, kpiItemId string, periodMonth, periodYear int) error {
	return r.DB.Where("person_id = ? AND kpi_item_id = ? AND period_month = ? AND period_year = ?", personId, kpiItemId, periodMonth, periodYear).
		Delete(&domainkpiitem.PersonKPITarget{}).Error
}

var _ interfacekpiitem.RepoKPIItemInterface = (*repo)(nil)
var _ interfacekpiitem.RepoPersonKPITargetInterface = (*personTargetRepo)(nil)
