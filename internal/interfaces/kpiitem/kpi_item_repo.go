package interfacekpiitem

import (
	domainkpiitem "teamleader-management/internal/domain/kpiitem"
	"teamleader-management/pkg/filter"
)

type RepoKPIItemInterface interface {
	Store(m domainkpiitem.KPIItem) error
	GetByID(id string) (domainkpiitem.KPIItem, error)
	GetByNameAndPillar(name, pillarId string) (domainkpiitem.KPIItem, error)
	GetAll(params filter.BaseParams) ([]domainkpiitem.KPIItem, int64, error)
	Update(m domainkpiitem.KPIItem) error
	Delete(id string) error
}

type RepoPersonKPITargetInterface interface {
	Upsert(target domainkpiitem.PersonKPITarget) error
	Get(personId, kpiItemId string, periodMonth, periodYear int) (domainkpiitem.PersonKPITarget, error)
	Delete(personId, kpiItemId string, periodMonth, periodYear int) error
}
