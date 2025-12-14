package interfacetltraining

import (
	domaintltraining "teamleader-management/internal/domain/tltraining"
	"teamleader-management/pkg/filter"
)

type RepoTLTrainingInterface interface {
	StoreMultiple(records []domaintltraining.TLTrainingParticipation) error
	GetByTrainingBatch(trainingBatch string) ([]domaintltraining.TLTrainingParticipation, error)
	GetAll(params filter.BaseParams) ([]domaintltraining.TLTrainingParticipation, int64, error)
}
