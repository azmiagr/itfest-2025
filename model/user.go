package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRegister struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
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

type CompetitionRegistrationRequest struct {
	FullName      string `json:"full_name" binding:"required"`
	StudentNumber string `json:"student_number" binding:"required"`
	University    string `json:"university" binding:"required"`
	Major         string `json:"major" binding:"required"`
}

type UpdateProfile struct {
	FullName      string `json:"full_name"`
	StudentNumber string `json:"student_number"`
	University    string `json:"university"`
	Major         string `json:"major"`
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
	OTP    string    `json:"otp" binding:"required"`
}

type ResetPasswordRequest struct {
	UserID          uuid.UUID `json:"user_id"`
	NewPassword     string    `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string    `json:"confirm_password" binding:"required,min=8"`
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

type GetUserPaymentStatus struct {
	FullName        string `json:"fullname"`
	StudentNumber   string `json:"student_number"`
	Email           string `json:"email"`
	PaymentTransc   string `json:"payment_transc"`
	TeamName        string `json:"team_name"`
	TeamStatus      string `json:"team_status"`
	CompetitionName string `json:"competition_name"`
}

type GetTotalParticipant struct {
	TotalUIUX int `json:"total_uiux"`
	TotalBP   int `json:"total_bp"`
}
