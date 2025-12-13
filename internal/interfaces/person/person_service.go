package interfaceperson

import (
	domainperson "teamleader-management/internal/domain/person"
	"teamleader-management/internal/dto"
	"teamleader-management/pkg/filter"
)

type ServicePersonInterface interface {
	Create(req dto.PersonCreate) (domainperson.Person, error)
	GetByID(id string) (domainperson.Person, error)
	GetAll(params filter.BaseParams) ([]domainperson.Person, int64, error)
	Update(id string, req dto.PersonUpdate) (domainperson.Person, error)
	Delete(id string) error
}
