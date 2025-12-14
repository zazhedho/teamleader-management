package repositorytlactivity

import (
	"fmt"
	domaintlactivity "teamleader-management/internal/domain/tlactivity"
	interfacetlactivity "teamleader-management/internal/interfaces/tlactivity"
	"teamleader-management/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewTLActivityRepo(db *gorm.DB) interfacetlactivity.RepoTLActivityInterface {
	return &repo{DB: db}
}

func (r *repo) Store(m domaintlactivity.TLDailyActivity) error {
	return r.DB.Create(&m).Error
}

func (r *repo) GetByID(id string) (domaintlactivity.TLDailyActivity, error) {
	var ret domaintlactivity.TLDailyActivity
	if err := r.DB.Where("id = ?", id).First(&ret).Error; err != nil {
		return domaintlactivity.TLDailyActivity{}, err
	}
	return ret, nil
}

func (r *repo) GetAll(params filter.BaseParams) ([]domaintlactivity.TLDailyActivity, int64, error) {
	var (
		ret       []domaintlactivity.TLDailyActivity
		totalData int64
	)

	query := r.DB.Model(&domaintlactivity.TLDailyActivity{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(activity_type) LIKE LOWER(?) OR LOWER(COALESCE(kecamatan, '')) LIKE LOWER(?) OR LOWER(COALESCE(desa, '')) LIKE LOWER(?)", searchPattern, searchPattern, searchPattern)
	}

	for key, value := range params.Filters {
		if value == nil {
			continue
		}

		switch key {
		case "person_id":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("person_id = ?", s)
			}
		case "activity_type":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("activity_type = ?", s)
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
			"person_id":     true,
			"date":          true,
			"activity_type": true,
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

func (r *repo) Update(m domaintlactivity.TLDailyActivity) error {
	return r.DB.Save(&m).Error
}

func (r *repo) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&domaintlactivity.TLDailyActivity{}).Error
}

var _ interfacetlactivity.RepoTLActivityInterface = (*repo)(nil)
