package servicekpiitem

import (
	"errors"
	"strings"
	"time"

	domainkpiitem "teamleader-management/internal/domain/kpiitem"
	"teamleader-management/internal/dto"
	interfacekpiitem "teamleader-management/internal/interfaces/kpiitem"
	interfaceperson "teamleader-management/internal/interfaces/person"
	interfacepillar "teamleader-management/internal/interfaces/pillar"
	"teamleader-management/pkg/filter"
	"teamleader-management/utils"
)

type ServiceKPIItem struct {
	KPIRepo          interfacekpiitem.RepoKPIItemInterface
	PillarRepo       interfacepillar.RepoPillarInterface
	PersonRepo       interfaceperson.RepoPersonInterface
	PersonTargetRepo interfacekpiitem.RepoPersonKPITargetInterface
}

func NewKPIItemService(kRepo interfacekpiitem.RepoKPIItemInterface, pRepo interfacepillar.RepoPillarInterface, personRepo interfaceperson.RepoPersonInterface, targetRepo interfacekpiitem.RepoPersonKPITargetInterface) *ServiceKPIItem {
	return &ServiceKPIItem{
		KPIRepo:          kRepo,
		PillarRepo:       pRepo,
		PersonRepo:       personRepo,
		PersonTargetRepo: targetRepo,
	}
}

func (s *ServiceKPIItem) Create(req dto.KPIItemCreate, actorId string) (domainkpiitem.KPIItem, error) {
	if _, err := s.PillarRepo.GetByID(req.PillarId); err != nil {
		return domainkpiitem.KPIItem{}, errors.New("pillar not found")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return domainkpiitem.KPIItem{}, errors.New("name is required")
	}

	if req.Weight < 0 || req.Weight > 100 {
		return domainkpiitem.KPIItem{}, errors.New("weight must be between 0 and 100")
	}

	if existing, err := s.KPIRepo.GetByNameAndPillar(name, req.PillarId); err == nil && existing.Id != "" {
		return domainkpiitem.KPIItem{}, errors.New("kpi item name already exists in this pillar")
	}

	entity := domainkpiitem.KPIItem{
		Id:                utils.CreateUUID(),
		PillarId:          req.PillarId,
		Name:              name,
		Weight:            req.Weight,
		TargetValue:       req.TargetValue,
		Unit:              req.Unit,
		Frequency:         req.Frequency,
		InputSource:       req.InputSource,
		AppliesToTL:       req.AppliesToTL,
		AppliesToSalesman: req.AppliesToSalesman,
		Notes:             req.Notes,
		CreatedAt:         time.Now(),
		CreatedBy:         actorId,
	}

	if err := s.KPIRepo.Store(entity); err != nil {
		return domainkpiitem.KPIItem{}, err
	}

	created, err := s.KPIRepo.GetByID(entity.Id)
	if err != nil {
		return domainkpiitem.KPIItem{}, err
	}
	return created, nil
}

func (s *ServiceKPIItem) GetByID(id string) (domainkpiitem.KPIItem, error) {
	return s.KPIRepo.GetByID(id)
}

func (s *ServiceKPIItem) GetAll(params filter.BaseParams) ([]domainkpiitem.KPIItem, int64, error) {
	params.Filters = filter.WhitelistFilter(params.Filters, []string{"pillar_id", "input_source", "applies_to_tl", "applies_to_salesman"})
	return s.KPIRepo.GetAll(params)
}

func (s *ServiceKPIItem) Update(id string, req dto.KPIItemUpdate, actorId string) (domainkpiitem.KPIItem, error) {
	kpi, err := s.KPIRepo.GetByID(id)
	if err != nil {
		return domainkpiitem.KPIItem{}, err
	}

	if req.PillarId != nil {
		if _, err := s.PillarRepo.GetByID(*req.PillarId); err != nil {
			return domainkpiitem.KPIItem{}, errors.New("pillar not found")
		}
		kpi.PillarId = *req.PillarId
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return domainkpiitem.KPIItem{}, errors.New("name cannot be empty")
		}
		if existing, err := s.KPIRepo.GetByNameAndPillar(name, kpi.PillarId); err == nil && existing.Id != "" && existing.Id != id {
			return domainkpiitem.KPIItem{}, errors.New("kpi item name already exists in this pillar")
		}
		kpi.Name = name
	}

	if req.Weight != nil {
		if *req.Weight < 0 || *req.Weight > 100 {
			return domainkpiitem.KPIItem{}, errors.New("weight must be between 0 and 100")
		}
		kpi.Weight = *req.Weight
	}

	if req.TargetValue != nil {
		kpi.TargetValue = req.TargetValue
	}

	if req.Unit != nil {
		kpi.Unit = req.Unit
	}

	if req.Frequency != nil {
		kpi.Frequency = req.Frequency
	}

	if req.InputSource != nil {
		kpi.InputSource = *req.InputSource
	}

	if req.AppliesToTL != nil {
		kpi.AppliesToTL = *req.AppliesToTL
	}

	if req.AppliesToSalesman != nil {
		kpi.AppliesToSalesman = *req.AppliesToSalesman
	}

	if req.Notes != nil {
		kpi.Notes = req.Notes
	}

	now := time.Now()
	kpi.UpdatedAt = now
	kpi.UpdatedBy = actorId

	if err := s.KPIRepo.Update(kpi); err != nil {
		return domainkpiitem.KPIItem{}, err
	}

	return kpi, nil
}

func (s *ServiceKPIItem) Delete(id string) error {
	_, err := s.KPIRepo.GetByID(id)
	if err != nil {
		return err
	}
	return s.KPIRepo.Delete(id)
}

func (s *ServiceKPIItem) UpsertPersonTarget(req dto.PersonKPITargetUpsert, actorId string) (domainkpiitem.PersonKPITarget, error) {
	if _, err := s.KPIRepo.GetByID(req.KPIItemId); err != nil {
		return domainkpiitem.PersonKPITarget{}, errors.New("kpi item not found")
	}
	if _, err := s.PersonRepo.GetByID(req.PersonId); err != nil {
		return domainkpiitem.PersonKPITarget{}, errors.New("person not found")
	}
	if req.PeriodMonth < 1 || req.PeriodMonth > 12 {
		return domainkpiitem.PersonKPITarget{}, errors.New("period_month must be 1-12")
	}
	if req.PeriodYear < 2000 {
		return domainkpiitem.PersonKPITarget{}, errors.New("period_year must be >= 2000")
	}

	existing, err := s.PersonTargetRepo.Get(req.PersonId, req.KPIItemId, req.PeriodMonth, req.PeriodYear)
	if err == nil && existing.Id != "" {
		now := time.Now()
		existing.TargetValue = req.TargetValue
		existing.UpdatedAt = now
		existing.UpdatedBy = actorId
		if err := s.PersonTargetRepo.Upsert(existing); err != nil {
			return domainkpiitem.PersonKPITarget{}, err
		}
		return s.PersonTargetRepo.Get(req.PersonId, req.KPIItemId, req.PeriodMonth, req.PeriodYear)
	}

	target := domainkpiitem.PersonKPITarget{
		Id:          utils.CreateUUID(),
		PersonId:    req.PersonId,
		KPIItemId:   req.KPIItemId,
		PeriodMonth: req.PeriodMonth,
		PeriodYear:  req.PeriodYear,
		TargetValue: req.TargetValue,
		CreatedAt:   time.Now(),
		CreatedBy:   actorId,
	}

	if err := s.PersonTargetRepo.Upsert(target); err != nil {
		return domainkpiitem.PersonKPITarget{}, err
	}

	return s.PersonTargetRepo.Get(req.PersonId, req.KPIItemId, req.PeriodMonth, req.PeriodYear)
}

func (s *ServiceKPIItem) DeletePersonTarget(personId, kpiItemId string, periodMonth, periodYear int, actorId string) error {
	return s.PersonTargetRepo.Delete(personId, kpiItemId, periodMonth, periodYear)
}

var _ interfacekpiitem.ServiceKPIItemInterface = (*ServiceKPIItem)(nil)
