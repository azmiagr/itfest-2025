package model

import "github.com/google/uuid"

type UserRegister struct {
	FullName         string `json:"full_Name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=8"`
	StudentNumber    string `json:"student_number" binding:"required"`
	TeamName         string `json:"team_name" binding:"required"`
	Major            string `json:"major" binding:"required"`
	University       string `json:"university" binding:"required"`
	RegistrationLink string `json:"gdrive_link" binding:"required"`
	PaymentTransc    string `json:"payment_transc"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
	RoleID int       `json:"role_id"`
}

type UserParam struct {
	UserID   uuid.UUID `json:"-"`
	Email    string    `json:"-"`
	Password string    `json:"-"`
	RoleID   int       `json:"-"`
}

type VerifyUser struct {
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	OtpCode string    `json:"otp_code" binding:"required"`
}
