package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID        uuid.UUID `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	Username      string    `json:"username" gorm:"type:varchar(70)"`
	Password      string    `json:"password" gorm:"type:varchar(80)"`
	Email         string    `json:"email" gorm:"type:varchar(50)"`
	PhoneNumber   string    `json:"phone_number" gorm:"type:varchar(15)"`
	GdriveLink    string    `json:"gdrive_link" gorm:"type:varchar(100)"`
	PaymentTransc string    `json:"payment_transc" gorm:"type:text"`
	RoleID        int       `json:"role_id"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Team          Team      `json:"team" gorm:"foreignKey:UserID"`
}
