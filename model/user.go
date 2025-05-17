package model

import (
	"time"

	"github.com/google/uuid"
)

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

type UserProfile struct {
	FullName      string `json:"full_name"`
	StudentNumber string `json:"student_number"`
	University    string `json:"university"`
	Major         string `json:"major"`
	Email         string `json:"email"`
}

type UpdateProfile struct {
	FullName      string `json:"full_name"`
	StudentNumber string `json:"student_number"`
	University    string `json:"university"`
	Major         string `json:"major"`
	Email         string `json:"email"`
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

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyToken struct {
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"token" binding:"required"`
}

type ResetPasswordRequest struct {
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type UserTeamProfile struct {
	LeaderName          string           `json:"leader_name"`
	StudentNumber       string           `json:"student_number"`
	CompetitionCategory string           `json:"competition_category"`
	Deadline            time.Time        `json:"deadline"`
	Members             []MemberResponse `json:"members"`
}

type MemberResponse struct {
	FullName      string `json:"full_name"`
	StudentNumber string `json:"student_number"`
}
