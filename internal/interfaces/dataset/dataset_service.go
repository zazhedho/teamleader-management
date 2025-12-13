package interfacedataset

import (
	"mime/multipart"
	domaindataset "teamleader-management/internal/domain/dataset"
	"teamleader-management/internal/dto"
)

type ServiceDatasetInterface interface {
	Create(datasetType string, req dto.DatasetUploadRequest, file multipart.File, fileHeader *multipart.FileHeader, actorId string) (domaindataset.DashboardDataset, []byte, error)
	List(filters map[string]interface{}) ([]domaindataset.DashboardDataset, int64, error)
	UpdateStatus(id string, status string, actorId string) (domaindataset.DashboardDataset, error)
}
