package interfacedataset

import (
	domaindataset "teamleader-management/internal/domain/dataset"
)

type RepoDatasetInterface interface {
	Store(m domaindataset.DashboardDataset) error
	GetByID(id string) (domaindataset.DashboardDataset, error)
	Update(m domaindataset.DashboardDataset) error
	GetAll(params map[string]interface{}) ([]domaindataset.DashboardDataset, int64, error)
}
