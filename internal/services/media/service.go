package servicemedia

import (
	"teamleader-management/utils"
	"time"

	domainmedia "teamleader-management/internal/domain/media"
	interfacemedia "teamleader-management/internal/interfaces/media"
)

type ServiceMedia struct {
	Repo interfacemedia.RepoMediaInterface
}

func NewMediaService(repo interfacemedia.RepoMediaInterface) *ServiceMedia {
	return &ServiceMedia{Repo: repo}
}

func (s *ServiceMedia) AttachMedia(entityType string, entityId string, fileUrls []string, fileNames []string, actorId string) ([]domainmedia.Media, error) {
	if len(fileUrls) == 0 {
		return []domainmedia.Media{}, nil
	}

	// Ensure fileNames has the same length as fileUrls
	if len(fileNames) < len(fileUrls) {
		for i := len(fileNames); i < len(fileUrls); i++ {
			fileNames = append(fileNames, "file")
		}
	}

	now := time.Now()
	var mediaRecords []domainmedia.Media

	for i, fileUrl := range fileUrls {
		media := domainmedia.Media{
			Id:           utils.CreateUUID(),
			EntityType:   entityType,
			EntityId:     entityId,
			FileUrl:      fileUrl,
			FileName:     fileNames[i],
			DisplayOrder: i + 1,
			CreatedAt:    now,
			CreatedBy:    actorId,
		}
		mediaRecords = append(mediaRecords, media)
	}

	if err := s.Repo.StoreMultiple(mediaRecords); err != nil {
		return nil, err
	}

	return mediaRecords, nil
}

func (s *ServiceMedia) GetMediaByEntity(entityType string, entityId string) ([]domainmedia.Media, error) {
	return s.Repo.GetByEntity(entityType, entityId)
}

func (s *ServiceMedia) DeleteMediaByEntity(entityType string, entityId string) error {
	return s.Repo.DeleteByEntity(entityType, entityId)
}

func (s *ServiceMedia) ReplaceMedia(entityType string, entityId string, fileUrls []string, fileNames []string, actorId string) ([]domainmedia.Media, error) {
	// Delete existing media
	if err := s.DeleteMediaByEntity(entityType, entityId); err != nil {
		return nil, err
	}

	// Attach new media
	return s.AttachMedia(entityType, entityId, fileUrls, fileNames, actorId)
}

var _ interfacemedia.ServiceMediaInterface = (*ServiceMedia)(nil)
