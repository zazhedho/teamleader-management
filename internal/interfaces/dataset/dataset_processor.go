package interfacedataset

import domaindataset "teamleader-management/internal/domain/dataset"

type DatasetProcessorInterface interface {
	ProcessStream(ds *domaindataset.DashboardDataset, data []byte, actorId string) (domaindataset.DashboardDataset, error)
}
