package service

import (
	"errors"
	"fmt"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/mail"
	"time"
)

type IOtpService interface {
	ResendOtp(param model.GetOtp) error
}

type OtpService struct {
	OtpRepository  repository.IOtpRepository
	UserRepository repository.IUserRepository
}

func NewOtpService(OtpRepository repository.IOtpRepository, UserRepository repository.IUserRepository) IOtpService {
	return &OtpService{
		OtpRepository:  OtpRepository,
		UserRepository: UserRepository,
	}
}

func (o *OtpService) ResendOtp(param model.GetOtp) error {
	user, err := o.UserRepository.GetUser(model.UserParam{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	if user.StatusAccount == "active" {
		return errors.New("your account is already active")
	}

	otp, err := o.OtpRepository.GetOtp(model.GetOtp{
		UserID: user.UserID,
	})
	if err != nil {
		return err
	}

	if otp.CreatedAt.After(time.Now().Add(-5 * time.Minute)) {
		return errors.New("you can only resend otp every 5 minutes")
	}

	otp.Code = mail.GenerateCode()

	err = mail.SendEmail(user.Email, "OTP Verification", fmt.Sprintf(`Your OTP code is %s`, otp.Code))
	if err != nil {
		return err
	}

	err = o.OtpRepository.UpdateOtp(otp)
	if err != nil {
		return err
	}

	return nil
}
