package repositorymedia

import (
	domainmedia "teamleader-management/internal/domain/media"
	interfacemedia "teamleader-management/internal/interfaces/media"

	"gorm.io/gorm"
)

type repo struct {
	DB *gorm.DB
}

func NewMediaRepo(db *gorm.DB) interfacemedia.RepoMediaInterface {
	return &repo{DB: db}
}

func (r *repo) Store(m domainmedia.Media) error {
	return r.DB.Create(&m).Error
}

func (r *repo) StoreMultiple(media []domainmedia.Media) error {
	return r.DB.Create(&media).Error
}

func (r *repo) GetByID(id string) (domainmedia.Media, error) {
	var ret domainmedia.Media
	if err := r.DB.Where("id = ?", id).First(&ret).Error; err != nil {
		return domainmedia.Media{}, err
	}
	return ret, nil
}

func (r *repo) GetByEntity(entityType string, entityId string) ([]domainmedia.Media, error) {
	var ret []domainmedia.Media
	if err := r.DB.Where("entity_type = ? AND entity_id = ?", entityType, entityId).
		Order("display_order ASC").
		Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *repo) DeleteByID(id string) error {
	return r.DB.Where("id = ?", id).Delete(&domainmedia.Media{}).Error
}

func (r *repo) DeleteByEntity(entityType string, entityId string) error {
	return r.DB.Where("entity_type = ? AND entity_id = ?", entityType, entityId).
		Delete(&domainmedia.Media{}).Error
}

var _ interfacemedia.RepoMediaInterface = (*repo)(nil)
