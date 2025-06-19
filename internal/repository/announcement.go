package repository

import (
	"itfest-2025/entity"

	"gorm.io/gorm"
)

type IAnnouncementRepository interface {
	CreateAnnouncement(tx *gorm.DB, req entity.Announcement) error
	GetAnnouncement() ([]*entity.Announcement, error)
}

type AnnouncementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) IAnnouncementRepository {
	return &AnnouncementRepository{
		db: db,
	}
}

func (r *AnnouncementRepository) CreateAnnouncement(tx *gorm.DB, req entity.Announcement) error {
	err := tx.Debug().Create(&req).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AnnouncementRepository) GetAnnouncement() ([]*entity.Announcement, error) {
	var announcement []*entity.Announcement
	err := r.db.Debug().Find(&announcement).Error
	if err != nil {
		return nil, err
	}

	return announcement, nil
}
