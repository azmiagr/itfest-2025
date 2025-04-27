package repository

import (
	"itfest-2025/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	GetUserByID(userID uuid.UUID) (*entity.User, error)
	UpdateUser(user *entity.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	err := u.db.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) GetUserByID(userID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := u.db.Debug().Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) UpdateUser(user *entity.User) error {
	err := u.db.Debug().Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}
