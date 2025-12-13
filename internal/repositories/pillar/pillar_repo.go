package repositorypillar

import (
	"fmt"

	domainpillar "teamleader-management/internal/domain/pillar"
	interfacepillar "teamleader-management/internal/interfaces/pillar"
	"teamleader-management/pkg/filter"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewPillarRepo(db *gorm.DB) interfacepillar.RepoPillarInterface {
	return &repo{DB: db}
}

func (r *repo) Store(m domainpillar.Pillar) error {
	return r.DB.Create(&m).Error
}

func (r *repo) GetByID(id string) (domainpillar.Pillar, error) {
	var ret domainpillar.Pillar
	if err := r.DB.Where("id = ?", id).First(&ret).Error; err != nil {
		return domainpillar.Pillar{}, err
	}
	return ret, nil
}

func (r *repo) GetByName(name string) (domainpillar.Pillar, error) {
	var ret domainpillar.Pillar
	if err := r.DB.Where("LOWER(name) = LOWER(?)", name).First(&ret).Error; err != nil {
		return domainpillar.Pillar{}, err
	}
	return ret, nil
}

func (r *repo) GetAll(params filter.BaseParams) ([]domainpillar.Pillar, int64, error) {
	var (
		ret       []domainpillar.Pillar
		totalData int64
	)

	query := r.DB.Model(&domainpillar.Pillar{})

	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?)", searchPattern)
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	if params.OrderBy != "" && params.OrderDirection != "" {
		validColumns := map[string]bool{
			"name":       true,
			"weight":     true,
			"created_at": true,
			"updated_at": true,
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

func (r *repo) Update(m domainpillar.Pillar) error {
	return r.DB.Save(&m).Error
}

func (r *repo) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&domainpillar.Pillar{}).Error
}

var _ interfacepillar.RepoPillarInterface = (*repo)(nil)
