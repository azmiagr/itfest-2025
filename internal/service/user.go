package service

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/bcrypt"
	"itfest-2025/pkg/jwt"
	"itfest-2025/pkg/mail"
	"itfest-2025/pkg/supabase"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(param *model.UserRegister) (string, error)
	Login(param model.UserLogin) (model.LoginResponse, error)
	UploadPayment(userID uuid.UUID, file *multipart.FileHeader) (string, error)
	VerifyUser(param model.VerifyUser) error
}

type UserService struct {
	UserRepository repository.IUserRepository
	TeamRepository repository.ITeamRepository
	OtpRepository  repository.IOtpRepository
	BCrypt         bcrypt.Interface
	JwtAuth        jwt.Interface
	Supabase       supabase.Interface
}

func NewUserService(userRepository repository.IUserRepository, teamRepository repository.ITeamRepository, otpRepository repository.IOtpRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface) IUserService {
	return &UserService{
		UserRepository: userRepository,
		TeamRepository: teamRepository,
		OtpRepository:  otpRepository,
		BCrypt:         bcrypt,
		JwtAuth:        jwtAuth,
		Supabase:       supabase,
	}
}

func (u *UserService) Register(param *model.UserRegister) (string, error) {
	err := u.TeamRepository.GetTeamByName(param.TeamName)
	if err == nil {
		return "", errors.New("team name already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	_, err = u.UserRepository.GetUser(model.UserParam{
		Email: param.Email,
	})

	if err == nil {
		return "", errors.New("email already registered")
	}

	hash, err := u.BCrypt.GenerateFromPassword(param.Password)
	if err != nil {
		return "", err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	user := &entity.User{
		UserID:        id,
		Username:      param.Username,
		Password:      hash,
		Email:         param.Email,
		PhoneNumber:   param.PhoneNumber,
		GdriveLink:    param.GdriveLink,
		StatusAccount: "inactive",
		PaymentTransc: "",
		RoleID:        2,
	}

	_, err = u.UserRepository.CreateUser(user)
	if err != nil {
		return "", err
	}

	code := mail.GenerateCode()
	otp := &entity.OtpCode{
		OtpID:  uuid.New(),
		UserID: user.UserID,
		Code:   code,
	}

	err = u.OtpRepository.CreateOtp(otp)
	if err != nil {
		return "", err
	}

	err = mail.SendEmail(user.Email, "OTP Verification", "Your OTP verification code is "+code+".")
	if err != nil {
		return "", err
	}

	teamID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	team := &entity.Team{
		TeamID:     teamID,
		TeamName:   param.TeamName,
		University: param.University,
		UserID:     id,
	}

	err = u.TeamRepository.CreateTeam(team)
	if err != nil {
		return "", err
	}

	return id.String(), nil

}

func (u *UserService) Login(param model.UserLogin) (model.LoginResponse, error) {
	var result model.LoginResponse

	user, err := u.UserRepository.GetUser(model.UserParam{
		Email: param.Email,
	})

	if err != nil {
		return result, err
	}

	err = u.BCrypt.CompareAndHashPassword(user.Password, param.Password)
	if err != nil {
		return result, err
	}

	token, err := u.JwtAuth.CreateJWTToken(user.UserID)
	if err != nil {
		return result, errors.New("failed to create token")
	}

	result.UserID = user.UserID
	result.Token = token
	result.RoleID = user.RoleID

	return result, nil
}

func (u *UserService) UploadPayment(userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return "", errors.New("user not found")
	}

	signedURL, err := u.Supabase.UploadFile(file)
	if err != nil {
		return "", err
	}

	user.PaymentTransc = signedURL

	err = u.UserRepository.UpdateUser(user)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func (u *UserService) VerifyUser(param model.VerifyUser) error {
	otp, err := u.OtpRepository.GetOtp(model.GetOtp{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	if otp.Code != param.OtpCode {
		return errors.New("invalid otp code")
	}

	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		return err
	}

	expiredThreshold := time.Now().Add(-time.Duration(expiredTime) * time.Minute)
	if otp.CreatedAt.Before(expiredThreshold) {
		return errors.New("otp expired")
	}

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	user.StatusAccount = "active"
	err = u.UserRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	err = u.OtpRepository.DeleteOtp(otp)
	if err != nil {
		return err
	}

	return nil
}
