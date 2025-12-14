package repositorymetric

import (
	domainmetric "teamleader-management/internal/domain/metric"
	interfacemetric "teamleader-management/internal/interfaces/metric"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "period_date"}, {Name: "person_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"dataset_id":  clause.Expr{SQL: "EXCLUDED.dataset_id"},
			"honda_id":    clause.Expr{SQL: "EXCLUDED.honda_id"},
			"dealer_code": clause.Expr{SQL: "EXCLUDED.dealer_code"},
			"score":       clause.Expr{SQL: "EXCLUDED.score"},
			"pass_status": clause.Expr{SQL: "EXCLUDED.pass_status"},
			"updated_at":  clause.Expr{SQL: "EXCLUDED.updated_at"},
			"updated_by":  clause.Expr{SQL: "EXCLUDED.updated_by"},
		}),
	}).Create(&entries).Error
}

func (r *repo) SaveAppleLogins(entries []domainmetric.AppleLogin) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "period_date"}, {Name: "person_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"dataset_id":   clause.Expr{SQL: "EXCLUDED.dataset_id"},
			"honda_id":     clause.Expr{SQL: "EXCLUDED.honda_id"},
			"dealer_code":  clause.Expr{SQL: "EXCLUDED.dealer_code"},
			"login_date":   clause.Expr{SQL: "EXCLUDED.login_date"},
			"morning_done": clause.Expr{SQL: "EXCLUDED.morning_done"},
			"evening_done": clause.Expr{SQL: "EXCLUDED.evening_done"},
			"updated_at":   clause.Expr{SQL: "EXCLUDED.updated_at"},
			"updated_by":   clause.Expr{SQL: "EXCLUDED.updated_by"},
		}),
	}).Create(&entries).Error
}

func (r *repo) SaveSalesFLP(entries []domainmetric.SalesFLP) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "period_date"}, {Name: "person_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"dataset_id":  clause.Expr{SQL: "EXCLUDED.dataset_id"},
			"honda_id":    clause.Expr{SQL: "EXCLUDED.honda_id"},
			"dealer_code": clause.Expr{SQL: "EXCLUDED.dealer_code"},
			"flp_amount":  clause.Expr{SQL: "EXCLUDED.flp_amount"},
			"updated_at":  clause.Expr{SQL: "EXCLUDED.updated_at"},
			"updated_by":  clause.Expr{SQL: "EXCLUDED.updated_by"},
		}),
	}).Create(&entries).Error
}

func (r *repo) SaveApplePoints(entries []domainmetric.ApplePoint) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "period_date"}, {Name: "person_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"dataset_id": clause.Expr{SQL: "EXCLUDED.dataset_id"},
			"honda_id":   clause.Expr{SQL: "EXCLUDED.honda_id"},
			"points":     clause.Expr{SQL: "EXCLUDED.points"},
			"updated_at": clause.Expr{SQL: "EXCLUDED.updated_at"},
			"updated_by": clause.Expr{SQL: "EXCLUDED.updated_by"},
		}),
	}).Create(&entries).Error
}

func (r *repo) SaveMyHeroPoints(entries []domainmetric.MyHeroPoint) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "period_date"}, {Name: "person_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"dataset_id":  clause.Expr{SQL: "EXCLUDED.dataset_id"},
			"honda_id":    clause.Expr{SQL: "EXCLUDED.honda_id"},
			"dealer_code": clause.Expr{SQL: "EXCLUDED.dealer_code"},
			"points":      clause.Expr{SQL: "EXCLUDED.points"},
			"updated_at":  clause.Expr{SQL: "EXCLUDED.updated_at"},
			"updated_by":  clause.Expr{SQL: "EXCLUDED.updated_by"},
		}),
	}).Create(&entries).Error
}

func (r *repo) SaveProspects(entries []domainmetric.Prospect) error {
	if len(entries) == 0 {
		return nil
	}
	return r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "period_date"}, {Name: "person_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"dataset_id":     clause.Expr{SQL: "EXCLUDED.dataset_id"},
			"honda_id":       clause.Expr{SQL: "EXCLUDED.honda_id"},
			"prospect_count": clause.Expr{SQL: "EXCLUDED.prospect_count"},
			"updated_at":     clause.Expr{SQL: "EXCLUDED.updated_at"},
			"updated_by":     clause.Expr{SQL: "EXCLUDED.updated_by"},
		}),
	}).Create(&entries).Error
}

var _ interfacemetric.RepoMetricInterface = (*repo)(nil)
