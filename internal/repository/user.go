package repository

import (
	"itfest-2025/entity"
	"itfest-2025/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(tx *gorm.DB, user *entity.User) (*entity.User, error)
	UpdateUser(tx *gorm.DB, user *entity.User) error
	GetUser(param model.UserParam) (*entity.User, error)
	GetAllUser() ([]*entity.User, error)
	GetCountPayment() (int64, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(tx *gorm.DB, user *entity.User) (*entity.User, error) {
	err := tx.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) GetUser(param model.UserParam) (*entity.User, error) {
	user := entity.User{}
	err := u.db.Debug().Preload("Team").Where(&param).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) UpdateUser(tx *gorm.DB, user *entity.User) error {
	err := tx.Where("user_id = ?", user.UserID).Updates(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) GetAllUser() ([]*entity.User, error) {
	var users []*entity.User
	err := u.db.Debug().Preload("Team.TeamMembers").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserRepository) GetCountPayment() (int64, error) {
	var count int64
	err := u.db.Debug().Model(&entity.User{}).Where("registration_link IS NOT NULL AND payment_transc <> ''").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
