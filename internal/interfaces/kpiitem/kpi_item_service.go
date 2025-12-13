package interfacekpiitem

import (
	domainkpiitem "teamleader-management/internal/domain/kpiitem"
	"teamleader-management/internal/dto"
	"teamleader-management/pkg/filter"
)

type ServiceKPIItemInterface interface {
	Create(req dto.KPIItemCreate, actorId string) (domainkpiitem.KPIItem, error)
	GetByID(id string) (domainkpiitem.KPIItem, error)
	GetAll(params filter.BaseParams) ([]domainkpiitem.KPIItem, int64, error)
	Update(id string, req dto.KPIItemUpdate, actorId string) (domainkpiitem.KPIItem, error)
	Delete(id string) error

	UpsertPersonTarget(req dto.PersonKPITargetUpsert, actorId string) (domainkpiitem.PersonKPITarget, error)
	DeletePersonTarget(personId, kpiItemId string, periodMonth, periodYear int, actorId string) error
}
