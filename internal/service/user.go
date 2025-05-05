package service

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/bcrypt"
	"itfest-2025/pkg/database/mariadb"
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
	GetUser(param model.UserParam) (*entity.User, error)
}

type UserService struct {
	db             *gorm.DB
	UserRepository repository.IUserRepository
	TeamRepository repository.ITeamRepository
	OtpRepository  repository.IOtpRepository
	BCrypt         bcrypt.Interface
	JwtAuth        jwt.Interface
	Supabase       supabase.Interface
}

func NewUserService(userRepository repository.IUserRepository, teamRepository repository.ITeamRepository, otpRepository repository.IOtpRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface) IUserService {
	return &UserService{
		db:             mariadb.Connection,
		UserRepository: userRepository,
		TeamRepository: teamRepository,
		OtpRepository:  otpRepository,
		BCrypt:         bcrypt,
		JwtAuth:        jwtAuth,
		Supabase:       supabase,
	}
}

func (u *UserService) Register(param *model.UserRegister) (string, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	err := u.TeamRepository.GetTeamByName(tx, param.TeamName)
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
		UserID:           id,
		FullName:         param.FullName,
		Email:            param.Email,
		Password:         hash,
		StudentNumber:    param.StudentNumber,
		RegistrationLink: param.RegistrationLink,
		StatusAccount:    "inactive",
		PaymentTransc:    "",
		RoleID:           2,
	}

	_, err = u.UserRepository.CreateUser(tx, user)
	if err != nil {
		return "", err
	}

	code := mail.GenerateCode()
	otp := &entity.OtpCode{
		OtpID:  uuid.New(),
		UserID: user.UserID,
		Code:   code,
	}

	err = u.OtpRepository.CreateOtp(tx, otp)
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
		Major:      param.Major,
		TeamStatus: "belum terverifikasi",
		UserID:     id,
	}

	err = u.TeamRepository.CreateTeam(tx, team)
	if err != nil {
		return "", err
	}

	err = tx.Commit().Error
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (u *UserService) Login(param model.UserLogin) (model.LoginResponse, error) {
	var result model.LoginResponse

	tx := u.db.Begin()
	defer tx.Rollback()

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

	err = tx.Commit().Error
	if err != nil {
		return result, nil
	}

	return result, nil
}

func (u *UserService) UploadPayment(userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

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

	err = u.UserRepository.UpdateUser(tx, user)
	if err != nil {
		return "", err
	}

	err = tx.Commit().Error
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func (u *UserService) VerifyUser(param model.VerifyUser) error {
	tx := u.db.Begin()
	defer tx.Rollback()

	otp, err := u.OtpRepository.GetOtp(tx, model.GetOtp{
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

	expiredThreshold := time.Now().UTC().Add(-time.Duration(expiredTime) * time.Minute)
	if otp.UpdatedAt.Before(expiredThreshold) {
		return errors.New("otp expired")
	}

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	user.StatusAccount = "active"
	err = u.UserRepository.UpdateUser(tx, user)
	if err != nil {
		return err
	}

	err = u.OtpRepository.DeleteOtp(tx, otp)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUser(param model.UserParam) (*entity.User, error) {
	return u.UserRepository.GetUser(param)
}
