package servicedataset

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	domaindataset "teamleader-management/internal/domain/dataset"
	"teamleader-management/internal/dto"
	interfacedataset "teamleader-management/internal/interfaces/dataset"
	"teamleader-management/utils"
)

type ServiceDataset struct {
	Repo interfacedataset.RepoDatasetInterface
}

func NewDatasetService(repo interfacedataset.RepoDatasetInterface) *ServiceDataset {
	return &ServiceDataset{Repo: repo}
}

func (s *ServiceDataset) Create(datasetType string, req dto.DatasetUploadRequest, file multipart.File, fileHeader *multipart.FileHeader, actorId string) (domaindataset.DashboardDataset, []byte, error) {
	dType := strings.ToUpper(strings.TrimSpace(datasetType))
	if !utils.AllowedDatasetTypes[dType] {
		return domaindataset.DashboardDataset{}, nil, fmt.Errorf("invalid dataset type: %s", dType)
	}

	if req.PeriodDate == "" {
		return domaindataset.DashboardDataset{}, nil, errors.New("period_date is required")
	}
	periodDate, err := time.Parse("2006-01-02", req.PeriodDate)
	if err != nil {
		return domaindataset.DashboardDataset{}, nil, errors.New("invalid period_date format, expected YYYY-MM-DD")
	}

	periodMonth := int(periodDate.Month())
	periodYear := periodDate.Year()
	if req.PeriodMonth >= 1 && req.PeriodMonth <= 12 {
		periodMonth = req.PeriodMonth
	}
	if req.PeriodYear >= 2000 {
		periodYear = req.PeriodYear
	}

	periodFrequency := strings.ToUpper(strings.TrimSpace(req.PeriodFrequency))
	if periodFrequency == "" {
		periodFrequency = "DAILY"
	}
	if _, ok := utils.AllowedPeriodFrequencies[periodFrequency]; !ok {
		return domaindataset.DashboardDataset{}, nil, errors.New("invalid period_frequency")
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return domaindataset.DashboardDataset{}, nil, fmt.Errorf("failed to read uploaded file: %w", err)
	}

	entity := domaindataset.DashboardDataset{
		Id:              utils.CreateUUID(),
		Type:            dType,
		PeriodDate:      periodDate,
		PeriodMonth:     periodMonth,
		PeriodYear:      periodYear,
		PeriodFrequency: periodFrequency,
		FileName:        filepath.Base(fileHeader.Filename),
		UploadedBy:      actorId,
		UploadedAt:      time.Now(),
		Status:          utils.DatasetStatusProcessing,
		CreatedAt:       time.Now(),
		CreatedBy:       actorId,
	}

	if err := s.Repo.Store(entity); err != nil {
		return domaindataset.DashboardDataset{}, nil, err
	}

	return entity, data, nil
}

func (s *ServiceDataset) List(filters map[string]interface{}) ([]domaindataset.DashboardDataset, int64, error) {
	return s.Repo.GetAll(filters)
}

func (s *ServiceDataset) UpdateStatus(id string, status string, actorId string) (domaindataset.DashboardDataset, error) {
	status = strings.ToUpper(strings.TrimSpace(status))
	if status == "" {
		return domaindataset.DashboardDataset{}, errors.New("status is required")
	}

	ds, err := s.Repo.GetByID(id)
	if err != nil {
		return domaindataset.DashboardDataset{}, err
	}

	now := time.Now()
	ds.Status = status
	ds.UpdatedAt = now
	ds.UpdatedBy = actorId

	if err := s.Repo.Update(ds); err != nil {
		return domaindataset.DashboardDataset{}, err
	}

	return ds, nil
}

var _ interfacedataset.ServiceDatasetInterface = (*ServiceDataset)(nil)
