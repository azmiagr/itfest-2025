package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID           uuid.UUID `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	FullName         string    `json:"full_name" gorm:"type:varchar(70);"`
	Password         string    `json:"password" gorm:"type:varchar(80);not null"`
	Email            string    `json:"email" gorm:"type:varchar(50);not null"`
	PhoneNumber      string    `json:"phone_number" gorm:"type:varchar(20);"`
	StudentNumber    string    `json:"student_number" gorm:"type:varchar(20);"`
	RegistrationLink string    `json:"registration_link" gorm:"type:varchar(100);"`
	PaymentTransc    string    `json:"payment_transc" gorm:"type:text"`
	StatusAccount    string    `json:"-" gorm:"type:enum('inactive', 'active');"`
	StudentCardLink  string    `json:"student_card_link" gorm:"type:text"`
	University       string    `json:"university" gorm:"type:varchar(80);"`
	Major            string    `json:"major" gorm:"type:varchar(80);"`
	RoleID           int       `json:"role_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Team    Team      `json:"team" gorm:"foreignKey:UserID"`
	OtpCode []OtpCode `json:"otp_code" gorm:"foreignKey:UserID"`
}
