package interfacetlattendance

import (
	domaintlattendance "teamleader-management/internal/domain/tlattendance"
	"teamleader-management/internal/dto"
	"teamleader-management/pkg/filter"
)

type ServiceTLAttendanceInterface interface {
	Create(personId string, req dto.TLAttendanceCreate, actorId string) ([]domaintlattendance.TLAttendanceRecord, error)
	GetByRecordUniqueId(recordUniqueId string, personId string) ([]domaintlattendance.TLAttendanceRecord, error)
	GetAll(personId string, params filter.BaseParams) ([]domaintlattendance.TLAttendanceRecord, int64, error)
	Update(recordUniqueId string, personId string, req dto.TLAttendanceUpdate, actorId string) ([]domaintlattendance.TLAttendanceRecord, error)
	Delete(recordUniqueId string, personId string) error
}
