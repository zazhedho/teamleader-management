package interfacetlsession

import (
	"context"

	domaintlsession "teamleader-management/internal/domain/tlsession"
	"teamleader-management/internal/dto"
	"teamleader-management/pkg/filter"
)

type ServiceTLSessionInterface interface {
	Create(personId string, req dto.TLSessionCreate, actorId string) (domaintlsession.TLSession, error)
	GetByID(id string, personId string) (domaintlsession.TLSession, error)
	GetAll(personId string, params filter.BaseParams) ([]domaintlsession.TLSession, int64, error)
	Update(id string, personId string, req dto.TLSessionUpdate, actorId string) (domaintlsession.TLSession, error)
	Delete(ctx context.Context, id string, personId string) error
}
