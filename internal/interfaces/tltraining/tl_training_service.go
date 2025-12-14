package interfacetltraining

import (
	domaintltraining "teamleader-management/internal/domain/tltraining"
	"teamleader-management/internal/dto"
	"teamleader-management/pkg/filter"
)

type ServiceTLTrainingInterface interface {
	Create(personId string, req dto.TLTrainingCreate, actorId string) ([]domaintltraining.TLTrainingParticipation, error)
	GetByTrainingBatch(trainingBatch string, personId string) ([]domaintltraining.TLTrainingParticipation, error)
	GetAll(personId string, params filter.BaseParams) ([]domaintltraining.TLTrainingParticipation, int64, error)
}
