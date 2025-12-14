package servicetlattendance

import (
	"errors"
	"fmt"
	"time"

	domaintlattendance "teamleader-management/internal/domain/tlattendance"
	"teamleader-management/internal/dto"
	interfacetlattendance "teamleader-management/internal/interfaces/tlattendance"
	"teamleader-management/pkg/filter"

	"github.com/google/uuid"
)

type ServiceTLAttendance struct {
	Repo interfacetlattendance.RepoTLAttendanceInterface
}

func NewTLAttendanceService(repo interfacetlattendance.RepoTLAttendanceInterface) *ServiceTLAttendance {
	return &ServiceTLAttendance{Repo: repo}
}

func (s *ServiceTLAttendance) Create(personId string, req dto.TLAttendanceCreate, actorId string) ([]domaintlattendance.TLAttendanceRecord, error) {
	recordUniqueId := uuid.New().String()
	now := time.Now()

	var records []domaintlattendance.TLAttendanceRecord
	for _, att := range req.Attendance {
		record := domaintlattendance.TLAttendanceRecord{
			Id:             uuid.New().String(),
			TlPersonId:     personId,
			SalesmanId:     att.SalesmanPersonId,
			SalesmanName:   att.SalesmanName,
			Date:           req.Date,
			Status:         att.Status,
			RecordUniqueId: recordUniqueId,
			CreatedAt:      now,
			CreatedBy:      actorId,
			UpdatedAt:      now,
			UpdatedBy:      actorId,
		}
		records = append(records, record)
	}

	if err := s.Repo.StoreMultiple(records); err != nil {
		return nil, err
	}

	created, err := s.Repo.GetByRecordUniqueId(recordUniqueId)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *ServiceTLAttendance) GetByRecordUniqueId(recordUniqueId string, personId string) ([]domaintlattendance.TLAttendanceRecord, error) {
	records, err := s.Repo.GetByRecordUniqueId(recordUniqueId)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("attendance record not found")
	}

	if records[0].TlPersonId != personId {
		return nil, errors.New("unauthorized access to this attendance record")
	}

	return records, nil
}

func (s *ServiceTLAttendance) GetAll(personId string, params filter.BaseParams) ([]domaintlattendance.TLAttendanceRecord, int64, error) {
	params.Filters = filter.WhitelistFilter(params.Filters, []string{"tl_person_id", "salesman_id", "status", "date_from", "date_to"})
	params.Filters["tl_person_id"] = personId

	return s.Repo.GetAll(params)
}

func (s *ServiceTLAttendance) Update(recordUniqueId string, personId string, req dto.TLAttendanceUpdate, actorId string) ([]domaintlattendance.TLAttendanceRecord, error) {
	existingRecords, err := s.Repo.GetByRecordUniqueId(recordUniqueId)
	if err != nil {
		return nil, err
	}

	if len(existingRecords) == 0 {
		return nil, fmt.Errorf("attendance record not found")
	}

	if existingRecords[0].TlPersonId != personId {
		return nil, errors.New("unauthorized access to this attendance record")
	}

	if err := s.Repo.DeleteByRecordUniqueId(recordUniqueId); err != nil {
		return nil, err
	}

	dateToUse := existingRecords[0].Date
	if req.Date != nil {
		dateToUse = *req.Date
	}

	attendanceToUse := req.Attendance
	if len(attendanceToUse) == 0 {
		for _, existing := range existingRecords {
			attendanceToUse = append(attendanceToUse, dto.AttendanceRecord{
				SalesmanPersonId: existing.SalesmanId,
				SalesmanName:     existing.SalesmanName,
				Status:           existing.Status,
			})
		}
	}

	now := time.Now()
	var newRecords []domaintlattendance.TLAttendanceRecord
	for _, att := range attendanceToUse {
		record := domaintlattendance.TLAttendanceRecord{
			Id:             uuid.New().String(),
			TlPersonId:     personId,
			SalesmanId:     att.SalesmanPersonId,
			SalesmanName:   att.SalesmanName,
			Date:           dateToUse,
			Status:         att.Status,
			RecordUniqueId: recordUniqueId,
			CreatedAt:      existingRecords[0].CreatedAt,
			CreatedBy:      existingRecords[0].CreatedBy,
			UpdatedAt:      now,
			UpdatedBy:      actorId,
		}
		newRecords = append(newRecords, record)
	}

	if err := s.Repo.StoreMultiple(newRecords); err != nil {
		return nil, err
	}

	updated, err := s.Repo.GetByRecordUniqueId(recordUniqueId)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *ServiceTLAttendance) Delete(recordUniqueId string, personId string) error {
	records, err := s.Repo.GetByRecordUniqueId(recordUniqueId)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("attendance record not found")
	}

	if records[0].TlPersonId != personId {
		return errors.New("unauthorized access to this attendance record")
	}

	return s.Repo.DeleteByRecordUniqueId(recordUniqueId)
}

var _ interfacetlattendance.ServiceTLAttendanceInterface = (*ServiceTLAttendance)(nil)
