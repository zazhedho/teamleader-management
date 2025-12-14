package repositorytlsession

import (
	"fmt"
	domaintlsession "teamleader-management/internal/domain/tlsession"
	interfacetlsession "teamleader-management/internal/interfaces/tlsession"
	"teamleader-management/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewTLSessionRepo(db *gorm.DB) interfacetlsession.RepoTLSessionInterface {
	return &repo{DB: db}
}

func (r *repo) Store(m domaintlsession.TLSession) error {
	return r.DB.Create(&m).Error
}

func (r *repo) GetByID(id string) (domaintlsession.TLSession, error) {
	var ret domaintlsession.TLSession
	if err := r.DB.Where("id = ?", id).First(&ret).Error; err != nil {
		return domaintlsession.TLSession{}, err
	}
	return ret, nil
}

func (r *repo) GetAll(params filter.BaseParams) ([]domaintlsession.TLSession, int64, error) {
	var (
		ret       []domaintlsession.TLSession
		totalData int64
	)

	query := r.DB.Model(&domaintlsession.TLSession{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(COALESCE(notes, '')) LIKE LOWER(?)", searchPattern)
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
		case "session_type":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("session_type = ?", s)
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
			"id":           true,
			"person_id":    true,
			"session_type": true,
			"date":         true,
			"created_at":   true,
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

func (r *repo) Update(m domaintlsession.TLSession) error {
	return r.DB.Save(&m).Error
}

func (r *repo) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&domaintlsession.TLSession{}).Error
}

var _ interfacetlsession.RepoTLSessionInterface = (*repo)(nil)
