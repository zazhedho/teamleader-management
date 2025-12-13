package interfaceperson

import (
	domainperson "teamleader-management/internal/domain/person"
	"teamleader-management/pkg/filter"
)

type RepoPersonInterface interface {
	Store(m domainperson.Person) error
	GetByID(id string) (domainperson.Person, error)
	GetByHondaID(hondaId string) (domainperson.Person, error)
	GetAll(params filter.BaseParams) ([]domainperson.Person, int64, error)
	Update(m domainperson.Person) error
	Deactivate(id string) error
}
