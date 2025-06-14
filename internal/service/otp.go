package service

import (
	"errors"
	"fmt"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"
	"itfest-2025/pkg/mail"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOtpService interface {
	ResendOtp(param model.GetOtp) error
	ResendToken(userID uuid.UUID) error
}

type OtpService struct {
	db             *gorm.DB
	OtpRepository  repository.IOtpRepository
	UserRepository repository.IUserRepository
}

func NewOtpService(OtpRepository repository.IOtpRepository, UserRepository repository.IUserRepository) IOtpService {
	return &OtpService{
		db:             mariadb.Connection,
		OtpRepository:  OtpRepository,
		UserRepository: UserRepository,
	}
}

func (o *OtpService) ResendOtp(param model.GetOtp) error {
	tx := o.db.Begin()
	defer tx.Rollback()

	user, err := o.UserRepository.GetUser(model.UserParam{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	if user.StatusAccount == "active" {
		return errors.New("your account is already active")
	}

	otp, err := o.OtpRepository.GetOtp(tx, model.GetOtp{
		UserID: user.UserID,
	})
	if err != nil {
		return err
	}

	if otp.UpdatedAt.After(time.Now().UTC().Add(-5 * time.Minute)) {
		return errors.New("you can only resend otp every 5 minutes")
	}

	otp.Code = mail.GenerateCode()

	err = mail.SendEmail(user.Email, "OTP Verification", fmt.Sprintf(`Your OTP code is %s`, otp.Code))
	if err != nil {
		return err
	}

	err = o.OtpRepository.UpdateOtp(tx, otp)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpService) ResendToken(userID uuid.UUID) error {
	tx := o.db.Begin()
	defer tx.Rollback()

	user, err := o.UserRepository.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return err
	}

	otp, err := o.OtpRepository.GetOtp(tx, model.GetOtp{
		UserID: user.UserID,
	})
	if err != nil {
		return err
	}

	if otp.UpdatedAt.After(time.Now().UTC().Add(-1 * time.Minute)) {
		return errors.New("you can only resend otp every 5 minutes")
	}

	otp.Code = mail.GenerateRandomString(6)

	err = mail.SendEmail(user.Email, "Reset Password Token", "Your Reset Password Code is "+otp.Code+".")
	if err != nil {
		return err
	}

	err = o.OtpRepository.UpdateOtp(tx, otp)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil

}
