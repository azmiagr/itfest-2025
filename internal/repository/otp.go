package repository

import (
	"itfest-2025/entity"
	"itfest-2025/model"

	"gorm.io/gorm"
)

type IOtpRepository interface {
	GetOtp(tx *gorm.DB, param model.GetOtp) (*entity.OtpCode, error)
	CreateOtp(tx *gorm.DB, otp *entity.OtpCode) error
	UpdateOtp(tx *gorm.DB, otp *entity.OtpCode) error
	DeleteOtp(tx *gorm.DB, otp *entity.OtpCode) error
}

type OtpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) IOtpRepository {
	return &OtpRepository{
		db: db,
	}
}

func (o *OtpRepository) GetOtp(tx *gorm.DB, param model.GetOtp) (*entity.OtpCode, error) {
	var otp *entity.OtpCode
	err := o.db.Debug().Where(&param).First(&otp).Error
	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (o *OtpRepository) CreateOtp(tx *gorm.DB, otp *entity.OtpCode) error {
	err := o.db.Debug().Create(otp).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpRepository) UpdateOtp(tx *gorm.DB, otp *entity.OtpCode) error {
	err := o.db.Debug().Updates(otp).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpRepository) DeleteOtp(tx *gorm.DB, otp *entity.OtpCode) error {
	err := o.db.Debug().Delete(otp).Error
	if err != nil {
		return err
	}

	return nil
}
