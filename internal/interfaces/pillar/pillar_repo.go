package interfacepillar

import (
	domainpillar "teamleader-management/internal/domain/pillar"
	"teamleader-management/pkg/filter"
)

type RepoPillarInterface interface {
	Store(m domainpillar.Pillar) error
	GetByID(id string) (domainpillar.Pillar, error)
	GetByName(name string) (domainpillar.Pillar, error)
	GetAll(params filter.BaseParams) ([]domainpillar.Pillar, int64, error)
	Update(m domainpillar.Pillar) error
	Delete(id string) error
}
