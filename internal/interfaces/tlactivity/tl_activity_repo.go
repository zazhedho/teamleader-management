package interfacetlactivity

import (
	domaintlactivity "teamleader-management/internal/domain/tlactivity"
	"teamleader-management/pkg/filter"
)

type RepoTLActivityInterface interface {
	Store(m domaintlactivity.TLDailyActivity) error
	GetByID(id string) (domaintlactivity.TLDailyActivity, error)
	GetAll(params filter.BaseParams) ([]domaintlactivity.TLDailyActivity, int64, error)
	Update(m domaintlactivity.TLDailyActivity) error
	Delete(id string) error
}
