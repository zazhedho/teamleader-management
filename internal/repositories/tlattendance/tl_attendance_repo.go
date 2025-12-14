package repositorytlattendance

import (
	"fmt"
	domaintlattendance "teamleader-management/internal/domain/tlattendance"
	interfacetlattendance "teamleader-management/internal/interfaces/tlattendance"
	"teamleader-management/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewTLAttendanceRepo(db *gorm.DB) interfacetlattendance.RepoTLAttendanceInterface {
	return &repo{DB: db}
}

func (r *repo) StoreMultiple(records []domaintlattendance.TLAttendanceRecord) error {
	return r.DB.Create(&records).Error
}

func (r *repo) GetByRecordUniqueId(recordUniqueId string) ([]domaintlattendance.TLAttendanceRecord, error) {
	var ret []domaintlattendance.TLAttendanceRecord
	if err := r.DB.Where("record_unique_id = ?", recordUniqueId).Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *repo) GetAll(params filter.BaseParams) ([]domaintlattendance.TLAttendanceRecord, int64, error) {
	var (
		ret       []domaintlattendance.TLAttendanceRecord
		totalData int64
	)

	query := r.DB.Model(&domaintlattendance.TLAttendanceRecord{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(salesman_name) LIKE LOWER(?)", searchPattern)
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
			"id":           true,
			"tl_person_id": true,
			"salesman_id":  true,
			"date":         true,
			"status":       true,
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

func (r *repo) DeleteByRecordUniqueId(recordUniqueId string) error {
	return r.DB.Where("record_unique_id = ?", recordUniqueId).Delete(&domaintlattendance.TLAttendanceRecord{}).Error
}

var _ interfacetlattendance.RepoTLAttendanceInterface = (*repo)(nil)
