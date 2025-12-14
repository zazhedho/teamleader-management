package interfacemedia

import (
	domainmedia "teamleader-management/internal/domain/media"
)

type RepoMediaInterface interface {
	Store(m domainmedia.Media) error
	StoreMultiple(media []domainmedia.Media) error
	GetByID(id string) (domainmedia.Media, error)
	GetByEntity(entityType string, entityId string) ([]domainmedia.Media, error)
	DeleteByID(id string) error
	DeleteByEntity(entityType string, entityId string) error
}
