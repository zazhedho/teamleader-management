package repositorymetric

import (
	domainmetric "teamleader-management/internal/domain/metric"
	interfacemetric "teamleader-management/internal/interfaces/metric"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewMetricRepo(db *gorm.DB) interfacemetric.RepoMetricInterface {
	return &repo{DB: db}
}

func (r *repo) SaveQuizResults(entries []domainmetric.QuizResult) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Create(&entries).Error
}

func (r *repo) SaveAppleLogins(entries []domainmetric.AppleLogin) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Create(&entries).Error
}

func (r *repo) SaveSalesFLP(entries []domainmetric.SalesFLP) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Create(&entries).Error
}

func (r *repo) SaveApplePoints(entries []domainmetric.ApplePoint) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Create(&entries).Error
}

func (r *repo) SaveMyHeroPoints(entries []domainmetric.MyHeroPoint) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Create(&entries).Error
}

func (r *repo) SaveProspects(entries []domainmetric.Prospect) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Create(&entries).Error
}

var _ interfacemetric.RepoMetricInterface = (*repo)(nil)
