package interfacetlattendance

import (
	domaintlattendance "teamleader-management/internal/domain/tlattendance"
	"teamleader-management/pkg/filter"
)

type RepoTLAttendanceInterface interface {
	StoreMultiple(records []domaintlattendance.TLAttendanceRecord) error
	GetByRecordUniqueId(recordUniqueId string) ([]domaintlattendance.TLAttendanceRecord, error)
	GetAll(params filter.BaseParams) ([]domaintlattendance.TLAttendanceRecord, int64, error)
	DeleteByRecordUniqueId(recordUniqueId string) error
}
