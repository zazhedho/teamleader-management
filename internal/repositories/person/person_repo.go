package repositoryperson

import (
	"fmt"
	domainperson "teamleader-management/internal/domain/person"
	interfaceperson "teamleader-management/internal/interfaces/person"
	"teamleader-management/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPersonRepo(db *gorm.DB) interfaceperson.RepoPersonInterface {
	return &repo{DB: db}
}

func (r *repo) Store(m domainperson.Person) error {
	return r.DB.Create(&m).Error
}

func (r *repo) GetByID(id string) (domainperson.Person, error) {
	var ret domainperson.Person
	if err := r.DB.Where("id = ?", id).First(&ret).Error; err != nil {
		return domainperson.Person{}, err
	}
	return ret, nil
}

func (r *repo) GetByHondaID(hondaId string) (domainperson.Person, error) {
	var ret domainperson.Person
	if err := r.DB.Where("honda_id = ?", hondaId).First(&ret).Error; err != nil {
		return domainperson.Person{}, err
	}
	return ret, nil
}

func (r *repo) GetAll(params filter.BaseParams) ([]domainperson.Person, int64, error) {
	var (
		ret       []domainperson.Person
		totalData int64
	)

	query := r.DB.Model(&domainperson.Person{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(honda_id) LIKE LOWER(?) OR LOWER(COALESCE(dealer_code, '')) LIKE LOWER(?)", searchPattern, searchPattern, searchPattern)
	}

	for key, value := range params.Filters {
		if value == nil {
			continue
		}

		switch key {
		case "role":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("role = ?", s)
			}
		case "active":
			switch v := value.(type) {
			case bool:
				query = query.Where("active = ?", v)
			case string:
				if v == "true" || v == "false" {
					query = query.Where("active = ?", v == "true")
				}
			}
		case "dealer_code":
			if s, ok := value.(string); ok && s != "" {
				query = query.Where("dealer_code = ?", s)
			}
		}
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	if params.OrderBy != "" && params.OrderDirection != "" {
		validColumns := map[string]bool{
			"id":          true,
			"honda_id":    true,
			"name":        true,
			"role":        true,
			"dealer_code": true,
			"active":      true,
			"created_at":  true,
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

func (r *repo) Update(m domainperson.Person) error {
	return r.DB.Save(&m).Error
}

func (r *repo) Deactivate(id string) error {
	return r.DB.Model(&domainperson.Person{}).Where("id = ?", id).Update("active", false).Error
}

var _ interfaceperson.RepoPersonInterface = (*repo)(nil)
