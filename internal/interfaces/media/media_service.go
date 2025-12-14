package interfacemedia

import (
	"context"
	"mime/multipart"

	domainmedia "teamleader-management/internal/domain/media"
)

type ServiceMediaInterface interface {
	GetMediaByEntity(entityType string, entityId string) ([]domainmedia.Media, error)
	DeleteMediaByEntity(ctx context.Context, entityType string, entityId string) error
	UploadAndAttach(ctx context.Context, entityType string, entityId string, file *multipart.FileHeader, actorId string) (domainmedia.Media, error)
	DeleteMediaByID(ctx context.Context, mediaId string) error
	GetMediaByID(mediaId string) (domainmedia.Media, error)
}
