package servicetltraining

import (
	"errors"
	"fmt"
	"time"

	domaintltraining "teamleader-management/internal/domain/tltraining"
	"teamleader-management/internal/dto"
	interfacetltraining "teamleader-management/internal/interfaces/tltraining"
	"teamleader-management/pkg/filter"

	"github.com/google/uuid"
)

type ServiceTLTraining struct {
	Repo interfacetltraining.RepoTLTrainingInterface
}

func NewTLTrainingService(repo interfacetltraining.RepoTLTrainingInterface) *ServiceTLTraining {
	return &ServiceTLTraining{Repo: repo}
}

func (s *ServiceTLTraining) Create(personId string, req dto.TLTrainingCreate, actorId string) ([]domaintltraining.TLTrainingParticipation, error) {
	trainingBatch := uuid.New().String()
	now := time.Now()

	var records []domaintltraining.TLTrainingParticipation
	for _, participant := range req.Participants {
		record := domaintltraining.TLTrainingParticipation{
			Id:            uuid.New().String(),
			TlPersonId:    personId,
			TrainingName:  req.TrainingName,
			Date:          req.Date,
			SalesmanId:    participant.SalesmanPersonId,
			SalesmanName:  participant.SalesmanName,
			Status:        participant.Status,
			TrainingBatch: trainingBatch,
			CreatedAt:     now,
			CreatedBy:     actorId,
			UpdatedAt:     now,
			UpdatedBy:     actorId,
		}
		records = append(records, record)
	}

	if err := s.Repo.StoreMultiple(records); err != nil {
		return nil, err
	}

	created, err := s.Repo.GetByTrainingBatch(trainingBatch)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *ServiceTLTraining) GetByTrainingBatch(trainingBatch string, personId string) ([]domaintltraining.TLTrainingParticipation, error) {
	records, err := s.Repo.GetByTrainingBatch(trainingBatch)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("training record not found")
	}

	if records[0].TlPersonId != personId {
		return nil, errors.New("unauthorized access to this training record")
	}

	return records, nil
}

func (s *ServiceTLTraining) GetAll(personId string, params filter.BaseParams) ([]domaintltraining.TLTrainingParticipation, int64, error) {
	params.Filters = filter.WhitelistFilter(params.Filters, []string{"tl_person_id", "salesman_id", "status", "date_from", "date_to"})
	params.Filters["tl_person_id"] = personId

	return s.Repo.GetAll(params)
}

var _ interfacetltraining.ServiceTLTrainingInterface = (*ServiceTLTraining)(nil)
