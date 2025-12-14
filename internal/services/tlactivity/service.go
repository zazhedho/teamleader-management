package servicetlactivity

import (
	"context"
	"errors"
	"teamleader-management/utils"
	"time"

	domaintlactivity "teamleader-management/internal/domain/tlactivity"
	"teamleader-management/internal/dto"
	interfacemedia "teamleader-management/internal/interfaces/media"
	interfacetlactivity "teamleader-management/internal/interfaces/tlactivity"
	"teamleader-management/pkg/filter"
)

type ServiceTLActivity struct {
	Repo         interfacetlactivity.RepoTLActivityInterface
	MediaService interfacemedia.ServiceMediaInterface
}

func NewTLActivityService(repo interfacetlactivity.RepoTLActivityInterface, mediaService interfacemedia.ServiceMediaInterface) *ServiceTLActivity {
	return &ServiceTLActivity{
		Repo:         repo,
		MediaService: mediaService,
	}
}

func (s *ServiceTLActivity) Create(personId string, req dto.TLActivityCreate, actorId string) (domaintlactivity.TLDailyActivity, error) {
	activityId := utils.CreateUUID()

	entity := domaintlactivity.TLDailyActivity{
		Id:               activityId,
		PersonId:         personId,
		Date:             req.Date,
		ActivityType:     req.ActivityType,
		Kecamatan:        req.Kecamatan,
		Desa:             req.Desa,
		GpsLat:           req.GpsLat,
		GpsLng:           req.GpsLng,
		DurationHours:    req.DurationHours,
		ProspectCount:    req.ProspectCount,
		DealCount:        req.DealCount,
		MotorkuDownloads: req.MotorkuDownloads,
		Notes:            req.Notes,
		CreatedAt:        time.Now(),
		CreatedBy:        actorId,
		UpdatedAt:        time.Now(),
		UpdatedBy:        actorId,
	}

	if err := s.Repo.Store(entity); err != nil {
		return domaintlactivity.TLDailyActivity{}, err
	}

	created, err := s.Repo.GetByID(activityId)
	if err != nil {
		return domaintlactivity.TLDailyActivity{}, err
	}

	return created, nil
}

func (s *ServiceTLActivity) GetByID(id string, personId string) (domaintlactivity.TLDailyActivity, error) {
	activity, err := s.Repo.GetByID(id)
	if err != nil {
		return domaintlactivity.TLDailyActivity{}, err
	}

	if activity.PersonId != personId {
		return domaintlactivity.TLDailyActivity{}, errors.New("unauthorized access to this activity")
	}

	return activity, nil
}

func (s *ServiceTLActivity) GetAll(personId string, params filter.BaseParams) ([]domaintlactivity.TLDailyActivity, int64, error) {
	params.Filters = filter.WhitelistFilter(params.Filters, []string{"person_id", "activity_type", "date_from", "date_to"})
	params.Filters["person_id"] = personId

	return s.Repo.GetAll(params)
}

func (s *ServiceTLActivity) Update(id string, personId string, req dto.TLActivityUpdate, actorId string) (domaintlactivity.TLDailyActivity, error) {
	activity, err := s.Repo.GetByID(id)
	if err != nil {
		return domaintlactivity.TLDailyActivity{}, err
	}

	if activity.PersonId != personId {
		return domaintlactivity.TLDailyActivity{}, errors.New("unauthorized access to this activity")
	}

	if req.Date != nil {
		activity.Date = *req.Date
	}

	if req.ActivityType != nil {
		activity.ActivityType = *req.ActivityType
	}

	if req.Kecamatan != nil {
		activity.Kecamatan = req.Kecamatan
	}

	if req.Desa != nil {
		activity.Desa = req.Desa
	}

	if req.GpsLat != nil {
		activity.GpsLat = req.GpsLat
	}

	if req.GpsLng != nil {
		activity.GpsLng = req.GpsLng
	}

	if req.DurationHours != nil {
		activity.DurationHours = req.DurationHours
	}

	if req.ProspectCount != nil {
		activity.ProspectCount = *req.ProspectCount
	}

	if req.DealCount != nil {
		activity.DealCount = *req.DealCount
	}

	if req.MotorkuDownloads != nil {
		activity.MotorkuDownloads = *req.MotorkuDownloads
	}

	if req.Notes != nil {
		activity.Notes = req.Notes
	}

	activity.UpdatedAt = time.Now()
	activity.UpdatedBy = actorId

	if err := s.Repo.Update(activity); err != nil {
		return domaintlactivity.TLDailyActivity{}, err
	}

	return activity, nil
}

func (s *ServiceTLActivity) Delete(ctx context.Context, id string, personId string) error {
	activity, err := s.Repo.GetByID(id)
	if err != nil {
		return err
	}

	if activity.PersonId != personId {
		return errors.New("unauthorized access to this activity")
	}

	// Delete associated media (including from storage)
	if err := s.MediaService.DeleteMediaByEntity(ctx, utils.EntityTLActivity, id); err != nil {
		// Log error but continue with deletion
	}

	return s.Repo.Delete(id)
}

var _ interfacetlactivity.ServiceTLActivityInterface = (*ServiceTLActivity)(nil)
