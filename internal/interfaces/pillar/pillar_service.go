package interfacepillar

import (
	domainpillar "teamleader-management/internal/domain/pillar"
	"teamleader-management/internal/dto"
	"teamleader-management/pkg/filter"
)

type ServicePillarInterface interface {
	Create(req dto.PillarCreate) (domainpillar.Pillar, error)
	GetByID(id string) (domainpillar.Pillar, error)
	GetAll(params filter.BaseParams) ([]domainpillar.Pillar, int64, error)
	Update(id string, req dto.PillarUpdate) (domainpillar.Pillar, error)
	Delete(id string) error
}
