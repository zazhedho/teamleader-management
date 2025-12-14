package interfacetlsession

import (
	domaintlsession "teamleader-management/internal/domain/tlsession"
	"teamleader-management/pkg/filter"
)

type RepoTLSessionInterface interface {
	Store(m domaintlsession.TLSession) error
	GetByID(id string) (domaintlsession.TLSession, error)
	GetAll(params filter.BaseParams) ([]domaintlsession.TLSession, int64, error)
	Update(m domaintlsession.TLSession) error
	Delete(id string) error
}
