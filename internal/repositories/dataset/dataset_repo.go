package repositorydataset

import (
	domaindataset "teamleader-management/internal/domain/dataset"
	interfacedataset "teamleader-management/internal/interfaces/dataset"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewDatasetRepo(db *gorm.DB) interfacedataset.RepoDatasetInterface {
	return &repo{DB: db}
}

func (r *repo) Store(m domaindataset.DashboardDataset) error {
	return r.DB.Create(&m).Error
}

func (r *repo) GetByID(id string) (domaindataset.DashboardDataset, error) {
	var ret domaindataset.DashboardDataset
	if err := r.DB.Where("id = ?", id).First(&ret).Error; err != nil {
		return domaindataset.DashboardDataset{}, err
	}
	return ret, nil
}

func (r *repo) Update(m domaindataset.DashboardDataset) error {
	return r.DB.Save(&m).Error
}

func (r *repo) GetAll(filters map[string]interface{}) ([]domaindataset.DashboardDataset, int64, error) {
	var (
		ret       []domaindataset.DashboardDataset
		totalData int64
	)

	query := r.DB.Model(&domaindataset.DashboardDataset{})

	if v, ok := filters["type"].(string); ok && v != "" {
		query = query.Where("type = ?", v)
	}
	if v, ok := filters["period_year"].(int); ok && v != 0 {
		query = query.Where("period_year = ?", v)
	}
	if v, ok := filters["period_month"].(int); ok && v != 0 {
		query = query.Where("period_month = ?", v)
	}
	if v, ok := filters["status"].(string); ok && v != "" {
		query = query.Where("status = ?", v)
	}

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("uploaded_at DESC").Find(&ret).Error; err != nil {
		return nil, 0, err
	}

	return ret, totalData, nil
}

var _ interfacedataset.RepoDatasetInterface = (*repo)(nil)
