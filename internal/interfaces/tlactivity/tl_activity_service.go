package interfacetlactivity

import (
	"context"

	domaintlactivity "teamleader-management/internal/domain/tlactivity"
	"teamleader-management/internal/dto"
	"teamleader-management/pkg/filter"
)

type ServiceTLActivityInterface interface {
	Create(personId string, req dto.TLActivityCreate, actorId string) (domaintlactivity.TLDailyActivity, error)
	GetByID(id string, personId string) (domaintlactivity.TLDailyActivity, error)
	GetAll(personId string, params filter.BaseParams) ([]domaintlactivity.TLDailyActivity, int64, error)
	Update(id string, personId string, req dto.TLActivityUpdate, actorId string) (domaintlactivity.TLDailyActivity, error)
	Delete(ctx context.Context, id string, personId string) error
}
