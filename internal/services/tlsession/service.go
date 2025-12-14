package servicetlsession

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	domaintlsession "teamleader-management/internal/domain/tlsession"
	"teamleader-management/internal/dto"
	interfacemedia "teamleader-management/internal/interfaces/media"
	interfacetlsession "teamleader-management/internal/interfaces/tlsession"
	"teamleader-management/pkg/filter"

	"github.com/google/uuid"
)

type ServiceTLSession struct {
	Repo         interfacetlsession.RepoTLSessionInterface
	MediaService interfacemedia.ServiceMediaInterface
}

func NewTLSessionService(repo interfacetlsession.RepoTLSessionInterface, mediaService interfacemedia.ServiceMediaInterface) *ServiceTLSession {
	return &ServiceTLSession{
		Repo:         repo,
		MediaService: mediaService,
	}
}

func (s *ServiceTLSession) Create(personId string, req dto.TLSessionCreate, actorId string) (domaintlsession.TLSession, error) {
	sessionId := uuid.New().String()

	entity := domaintlsession.TLSession{
		Id:            sessionId,
		PersonId:      personId,
		SessionType:   req.SessionType,
		Date:          req.Date,
		Notes:         req.Notes,
		Attendees:     pq.StringArray(req.Attendees),
		DurationHours: req.DurationHours,
		CreatedAt:     time.Now(),
		CreatedBy:     actorId,
		UpdatedAt:     time.Now(),
		UpdatedBy:     actorId,
	}

	if err := s.Repo.Store(entity); err != nil {
		return domaintlsession.TLSession{}, err
	}

	created, err := s.Repo.GetByID(sessionId)
	if err != nil {
		return domaintlsession.TLSession{}, err
	}

	return created, nil
}

func (s *ServiceTLSession) GetByID(id string, personId string) (domaintlsession.TLSession, error) {
	session, err := s.Repo.GetByID(id)
	if err != nil {
		return domaintlsession.TLSession{}, err
	}

	if session.PersonId != personId {
		return domaintlsession.TLSession{}, errors.New("unauthorized access to this session")
	}

	return session, nil
}

func (s *ServiceTLSession) GetAll(personId string, params filter.BaseParams) ([]domaintlsession.TLSession, int64, error) {
	params.Filters = filter.WhitelistFilter(params.Filters, []string{"person_id", "session_type", "date_from", "date_to"})
	params.Filters["person_id"] = personId

	return s.Repo.GetAll(params)
}

func (s *ServiceTLSession) Update(id string, personId string, req dto.TLSessionUpdate, actorId string) (domaintlsession.TLSession, error) {
	session, err := s.Repo.GetByID(id)
	if err != nil {
		return domaintlsession.TLSession{}, err
	}

	if session.PersonId != personId {
		return domaintlsession.TLSession{}, errors.New("unauthorized access to this session")
	}

	if req.SessionType != nil {
		session.SessionType = *req.SessionType
	}

	if req.Date != nil {
		session.Date = *req.Date
	}

	if req.Notes != nil {
		session.Notes = req.Notes
	}

	if req.Attendees != nil {
		session.Attendees = pq.StringArray(req.Attendees)
	}

	if req.DurationHours != nil {
		session.DurationHours = req.DurationHours
	}

	session.UpdatedAt = time.Now()
	session.UpdatedBy = actorId

	if err := s.Repo.Update(session); err != nil {
		return domaintlsession.TLSession{}, err
	}

	return session, nil
}

func (s *ServiceTLSession) Delete(ctx context.Context, id string, personId string) error {
	session, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	if session.PersonId != personId {
		return errors.New("unauthorized access to this session")
	}

	// Delete associated media (including from storage)
	if err := s.MediaService.DeleteMediaByEntity(ctx, "tl_session", id); err != nil {
		// Log error but continue with deletion
	}

	return s.Repo.Delete(id)
}

var _ interfacetlsession.ServiceTLSessionInterface = (*ServiceTLSession)(nil)
