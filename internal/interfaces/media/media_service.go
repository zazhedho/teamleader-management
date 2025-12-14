package interfacemedia

import (
	domainmedia "teamleader-management/internal/domain/media"
)

type ServiceMediaInterface interface {
	AttachMedia(entityType string, entityId string, fileUrls []string, fileNames []string, actorId string) ([]domainmedia.Media, error)
	GetMediaByEntity(entityType string, entityId string) ([]domainmedia.Media, error)
	DeleteMediaByEntity(entityType string, entityId string) error
	ReplaceMedia(entityType string, entityId string, fileUrls []string, fileNames []string, actorId string) ([]domainmedia.Media, error)
}
