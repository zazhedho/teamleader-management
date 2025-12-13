package servicepillar

import (
	"errors"
	"strings"

	domainpillar "teamleader-management/internal/domain/pillar"
	"teamleader-management/internal/dto"
	interfacepillar "teamleader-management/internal/interfaces/pillar"
	"teamleader-management/pkg/filter"
	"teamleader-management/utils"
)

type ServicePillar struct {
	Repo interfacepillar.RepoPillarInterface
}

func NewPillarService(repo interfacepillar.RepoPillarInterface) *ServicePillar {
	return &ServicePillar{Repo: repo}
}

func (s *ServicePillar) Create(req dto.PillarCreate) (domainpillar.Pillar, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return domainpillar.Pillar{}, errors.New("name is required")
	}

	if req.Weight < 0 || req.Weight > 100 {
		return domainpillar.Pillar{}, errors.New("weight must be between 0 and 100")
	}

	if existing, err := s.Repo.GetByName(name); err == nil && existing.Id != "" {
		return domainpillar.Pillar{}, errors.New("pillar name already exists")
	}

	entity := domainpillar.Pillar{
		Id:          utils.CreateUUID(),
		Name:        name,
		Description: req.Description,
		Weight:      req.Weight,
	}

	if err := s.Repo.Store(entity); err != nil {
		return domainpillar.Pillar{}, err
	}

	created, err := s.Repo.GetByID(entity.Id)
	if err != nil {
		return domainpillar.Pillar{}, err
	}
	return created, nil
}

func (s *ServicePillar) GetByID(id string) (domainpillar.Pillar, error) {
	return s.Repo.GetByID(id)
}

func (s *ServicePillar) GetAll(params filter.BaseParams) ([]domainpillar.Pillar, int64, error) {
	return s.Repo.GetAll(params)
}

func (s *ServicePillar) Update(id string, req dto.PillarUpdate) (domainpillar.Pillar, error) {
	pillar, err := s.Repo.GetByID(id)
	if err != nil {
		return domainpillar.Pillar{}, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return domainpillar.Pillar{}, errors.New("name cannot be empty")
		}
		if existing, err := s.Repo.GetByName(name); err == nil && existing.Id != "" && existing.Id != id {
			return domainpillar.Pillar{}, errors.New("pillar name already exists")
		}
		pillar.Name = name
	}

	if req.Description != nil {
		pillar.Description = req.Description
	}

	if req.Weight != nil {
		if *req.Weight < 0 || *req.Weight > 100 {
			return domainpillar.Pillar{}, errors.New("weight must be between 0 and 100")
		}
		pillar.Weight = *req.Weight
	}

	if err := s.Repo.Update(pillar); err != nil {
		return domainpillar.Pillar{}, err
	}

	return pillar, nil
}

func (s *ServicePillar) Delete(id string) error {
	_, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.Repo.Delete(id)
}

var _ interfacepillar.ServicePillarInterface = (*ServicePillar)(nil)
