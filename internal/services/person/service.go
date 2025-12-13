package serviceperson

import (
	"errors"
	domainperson "teamleader-management/internal/domain/person"
	"teamleader-management/internal/dto"
	interfaceperson "teamleader-management/internal/interfaces/person"
	"teamleader-management/pkg/filter"
	"teamleader-management/utils"
	"time"
)

type ServicePerson struct {
	Repo interfaceperson.RepoPersonInterface
}

func NewPersonService(repo interfaceperson.RepoPersonInterface) *ServicePerson {
	return &ServicePerson{Repo: repo}
}

func (s *ServicePerson) Create(req dto.PersonCreate, actorId string) (domainperson.Person, error) {
	if err := utils.ValidateRole(req.Role); err != nil {
		return domainperson.Person{}, err
	}

	if existing, err := s.Repo.GetByHondaID(req.HondaId); err == nil && existing.Id != "" {
		return domainperson.Person{}, errors.New("honda_id already exists")
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	entity := domainperson.Person{
		HondaId:    req.HondaId,
		Name:       req.Name,
		JobTitle:   req.JobTitle,
		Role:       req.Role,
		DealerCode: req.DealerCode,
		Active:     active,
		CreatedAt:  time.Now(),
		CreatedBy:  actorId,
	}

	if err := s.Repo.Store(entity); err != nil {
		return domainperson.Person{}, err
	}

	created, err := s.Repo.GetByHondaID(req.HondaId)
	if err != nil {
		return domainperson.Person{}, err
	}

	return created, nil
}

func (s *ServicePerson) GetByID(id string) (domainperson.Person, error) {
	return s.Repo.GetByID(id)
}

func (s *ServicePerson) GetAll(params filter.BaseParams) ([]domainperson.Person, int64, error) {
	params.Filters = filter.WhitelistFilter(params.Filters, []string{"role", "active", "dealer_code"})
	return s.Repo.GetAll(params)
}

func (s *ServicePerson) Update(id string, req dto.PersonUpdate, actorId string) (domainperson.Person, error) {
	person, err := s.Repo.GetByID(id)
	if err != nil {
		return domainperson.Person{}, err
	}

	if req.Role != nil {
		if err := utils.ValidateRole(*req.Role); err != nil {
			return domainperson.Person{}, err
		}
		person.Role = *req.Role
	}

	if req.HondaId != nil {
		if existing, err := s.Repo.GetByHondaID(*req.HondaId); err == nil && existing.Id != "" && existing.Id != id {
			return domainperson.Person{}, errors.New("honda_id already exists")
		}
		person.HondaId = *req.HondaId
	}

	if req.Name != nil {
		person.Name = *req.Name
	}

	if req.JobTitle != nil {
		person.JobTitle = req.JobTitle
	}

	if req.DealerCode != nil {
		person.DealerCode = req.DealerCode
	}

	if req.Active != nil {
		person.Active = *req.Active
	}

	now := time.Now()
	person.UpdatedAt = now
	person.UpdatedBy = actorId

	if err := s.Repo.Update(person); err != nil {
		return domainperson.Person{}, err
	}

	return person, nil
}

func (s *ServicePerson) Delete(id string) error {
	_, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.Repo.Deactivate(id)
}

var _ interfaceperson.ServicePersonInterface = (*ServicePerson)(nil)
