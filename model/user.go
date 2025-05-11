package model

import "github.com/google/uuid"

type UserRegister struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
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
