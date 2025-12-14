package servicemedia

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"teamleader-management/pkg/storage"
	"teamleader-management/utils"
	"time"

	domainmedia "teamleader-management/internal/domain/media"
	interfacemedia "teamleader-management/internal/interfaces/media"
)

type ServiceMedia struct {
	Repo    interfacemedia.RepoMediaInterface
	Storage storage.StorageProvider
}

func NewMediaService(repo interfacemedia.RepoMediaInterface, storageProvider storage.StorageProvider) *ServiceMedia {
	return &ServiceMedia{
		Repo:    repo,
		Storage: storageProvider,
	}
}

func (s *ServiceMedia) GetMediaByEntity(entityType string, entityId string) ([]domainmedia.Media, error) {
	return s.Repo.GetByEntity(entityType, entityId)
}

func (s *ServiceMedia) DeleteMediaByEntity(ctx context.Context, entityType string, entityId string) error {
	// Get all media for this entity first
	mediaList, err := s.Repo.GetByEntity(entityType, entityId)
	if err != nil {
		return fmt.Errorf("failed to get media list: %w", err)
	}

	// Delete files from storage
	if s.Storage != nil {
		for _, media := range mediaList {
			if media.FileUrl != "" {
				if err := s.Storage.DeleteFile(ctx, media.FileUrl); err != nil {
					fmt.Printf("warning: failed to delete file from storage: %v\n", err)
				}
			}
		}
	}

	return s.Repo.DeleteByEntity(entityType, entityId)
}

func (s *ServiceMedia) UploadAndAttach(ctx context.Context, entityType string, entityId string, file *multipart.FileHeader, actorId string) (domainmedia.Media, error) {
	if s.Storage == nil {
		return domainmedia.Media{}, fmt.Errorf("storage provider not configured")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return domainmedia.Media{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	folder := utils.UnderscoreToDash(entityType)

	// Upload to storage
	fileUrl, err := s.Storage.UploadFile(ctx, src, file, folder)
	if err != nil {
		return domainmedia.Media{}, fmt.Errorf("failed to upload file: %w", err)
	}

	// Get existing media count for display order
	existingMedia, _ := s.Repo.GetByEntity(entityType, entityId)
	displayOrder := len(existingMedia) + 1

	// Determine file type from content type
	contentType := file.Header.Get("Content-Type")
	fileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))

	media := domainmedia.Media{
		Id:           utils.CreateUUID(),
		EntityType:   entityType,
		EntityId:     entityId,
		FileUrl:      fileUrl,
		FileName:     fileName,
		FileType:     &contentType,
		FileSize:     &file.Size,
		DisplayOrder: displayOrder,
		CreatedAt:    time.Now(),
		CreatedBy:    actorId,
	}

	if err := s.Repo.Store(media); err != nil {
		// Try to delete uploaded file if db save fails
		_ = s.Storage.DeleteFile(ctx, fileUrl)
		return domainmedia.Media{}, fmt.Errorf("failed to save media record: %w", err)
	}

	return media, nil
}

func (s *ServiceMedia) GetMediaByID(mediaId string) (domainmedia.Media, error) {
	return s.Repo.GetByID(mediaId)
}

func (s *ServiceMedia) DeleteMediaByID(ctx context.Context, mediaId string) error {
	media, err := s.Repo.GetByID(mediaId)
	if err != nil {
		return fmt.Errorf("media not found: %w", err)
	}

	// Delete from storage if storage provider is configured
	if s.Storage != nil && media.FileUrl != "" {
		if err := s.Storage.DeleteFile(ctx, media.FileUrl); err != nil {
			// Log error but continue with DB deletion
			fmt.Printf("warning: failed to delete file from storage: %v\n", err)
		}
	}

	return s.Repo.DeleteByID(mediaId)
}

var _ interfacemedia.ServiceMediaInterface = (*ServiceMedia)(nil)
