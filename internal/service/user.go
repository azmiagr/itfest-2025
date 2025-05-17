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
	Register(param *model.UserRegister) (model.RegisterResponse, error)
	Login(param model.UserLogin) (model.LoginResponse, error)
	UploadPayment(userID uuid.UUID, file *multipart.FileHeader) (string, error)
	VerifyUser(param model.VerifyUser) error
	UpdateProfile(userID uuid.UUID, param model.UpdateProfile) error
	GetUserProfile(userID uuid.UUID) (model.UserProfile, error)
	GetMyTeamProfile(userID uuid.UUID) (*model.UserTeamProfile, error)
	ForgotPassword(email string) error
	ChangePasswordAfterVerify(userID uuid.UUID, param model.ResetPasswordRequest) error
	VerifyToken(param model.VerifyToken) error
	GetUser(param model.UserParam) (*entity.User, error)
}

type UserService struct {
	db                    *gorm.DB
	UserRepository        repository.IUserRepository
	TeamRepository        repository.ITeamRepository
	OtpRepository         repository.IOtpRepository
	CompetitionRepository repository.ICompetitionRepository
	BCrypt                bcrypt.Interface
	JwtAuth               jwt.Interface
	Supabase              supabase.Interface
}

func NewUserService(userRepository repository.IUserRepository, teamRepository repository.ITeamRepository, otpRepository repository.IOtpRepository, competitionRepository repository.ICompetitionRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface) IUserService {
	return &UserService{
		db:                    mariadb.Connection,
		UserRepository:        userRepository,
		TeamRepository:        teamRepository,
		OtpRepository:         otpRepository,
		CompetitionRepository: competitionRepository,
		BCrypt:                bcrypt,
		JwtAuth:               jwtAuth,
		Supabase:              supabase,
	}
}

func (u *UserService) Register(param *model.UserRegister) (model.RegisterResponse, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	var result model.RegisterResponse

	_, err := u.UserRepository.GetUser(model.UserParam{
		Email: param.Email,
	})

	if err == nil {
		return result, errors.New("email already registered")
	}

	hash, err := u.BCrypt.GenerateFromPassword(param.Password)
	if err != nil {
		return result, err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return result, err
	}

	user := &entity.User{
		UserID:        id,
		Email:         param.Email,
		Password:      hash,
		StatusAccount: "inactive",
		RoleID:        2,
	}

	_, err = u.UserRepository.CreateUser(tx, user)
	if err != nil {
		return result, err
	}

	token, err := u.JwtAuth.CreateJWTToken(user.UserID)
	if err != nil {
		return result, errors.New("failed to create token")
	}

	code := mail.GenerateCode()
	otp := &entity.OtpCode{
		OtpID:  uuid.New(),
		UserID: user.UserID,
		Code:   code,
	}

	err = u.OtpRepository.CreateOtp(tx, otp)
	if err != nil {
		return result, err
	}

	err = mail.SendEmail(user.Email, "OTP Verification", "Your OTP verification code is "+code+".")
	if err != nil {
		return result, err
	}

	err = tx.Commit().Error
	if err != nil {
		return result, err
	}

	result.Token = token

	return result, nil
}

func (u *UserService) Login(param model.UserLogin) (model.LoginResponse, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	var result model.LoginResponse

	user, err := u.UserRepository.GetUser(model.UserParam{
		Email: param.Email,
	})

	if err != nil {
		return result, err
	}

	err = u.BCrypt.CompareAndHashPassword(user.Password, param.Password)
	if err != nil {
		return result, errors.New("email or password is wrong")
	}

	token, err := u.JwtAuth.CreateJWTToken(user.UserID)
	if err != nil {
		return result, errors.New("failed to create token")
	}

	result.Token = token

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

	paymentURL, err := u.Supabase.UploadFile(file)
	if err != nil {
		return "", err
	}

	user.PaymentTransc = paymentURL

	err = u.UserRepository.UpdateUser(tx, user)
	if err != nil {
		return "", err
	}

	err = tx.Commit().Error
	if err != nil {
		return "", err
	}

	return paymentURL, nil
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

func (u *UserService) UpdateProfile(userID uuid.UUID, param model.UpdateProfile) error {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: userID,
	})

	if err != nil {
		return err
	}

	user.FullName = param.FullName
	user.StudentNumber = param.StudentNumber
	user.University = param.University
	user.Major = param.Major
	user.Email = param.Email

	err = u.UserRepository.UpdateUser(tx, user)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUserProfile(userID uuid.UUID) (model.UserProfile, error) {
	var result model.UserProfile

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return result, err
	}

	result.FullName = user.FullName
	result.StudentNumber = user.StudentNumber
	result.University = user.University
	result.Major = user.Major
	result.Email = user.Email

	return result, nil
}

func (u *UserService) GetUser(param model.UserParam) (*entity.User, error) {
	return u.UserRepository.GetUser(param)
}

func (u *UserService) GetMyTeamProfile(userID uuid.UUID) (*model.UserTeamProfile, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	team, err := u.TeamRepository.GetTeamByUserID(tx, userID)
	if err != nil {
		return nil, err
	}

	members, err := u.TeamRepository.GetTeamMemberByTeamID(tx, team.TeamID)
	if err != nil {
		return nil, err
	}

	var memberResponse []model.MemberResponse
	for _, v := range members {
		memberResponse = append(memberResponse, model.MemberResponse{
			FullName:      v.MemberName,
			StudentNumber: v.StudentNumber,
		})
	}

	competititon, err := u.CompetitionRepository.GetCompetitionByID(tx, team.CompetitionID)
	if err != nil {
		return nil, err
	}

	TeamProfileResponse := &model.UserTeamProfile{
		LeaderName:          user.FullName,
		StudentNumber:       user.StudentNumber,
		CompetitionCategory: competititon.CompetitionName,
		Members:             memberResponse,
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return TeamProfileResponse, nil

}

func (u *UserService) ForgotPassword(email string) error {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.UserRepository.GetUser(model.UserParam{
		Email: email,
	})
	if err != nil {
		return err
	}

	token := mail.GenerateRandomString(6)
	err = u.OtpRepository.CreateOtp(tx, &entity.OtpCode{
		OtpID:  uuid.New(),
		UserID: user.UserID,
		Code:   token,
	})
	if err != nil {
		return err
	}

	err = mail.SendEmail(user.Email, "Reset Password Token", "Your Reset Password Code is "+token+".")
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) VerifyToken(param model.VerifyToken) error {
	tx := u.db.Begin()
	defer tx.Rollback()

	otp, err := u.OtpRepository.GetOtp(tx, model.GetOtp{
		UserID: param.UserID,
	})
	if err != nil {
		return err
	}

	if otp.Code != param.Token {
		return errors.New("invalid token")
	}

	expiredTime, err := strconv.Atoi(os.Getenv("EXPIRED_OTP"))
	if err != nil {
		return err
	}

	expiredThreshold := time.Now().UTC().Add(-time.Duration(expiredTime) * time.Minute)
	if otp.UpdatedAt.Before(expiredThreshold) {
		return errors.New("token expired")
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

func (u *UserService) ChangePasswordAfterVerify(userID uuid.UUID, param model.ResetPasswordRequest) error {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return err
	}

	if param.NewPassword != param.ConfirmPassword {
		return errors.New("password mismatch")
	}

	hashPassword, err := u.BCrypt.GenerateFromPassword(param.NewPassword)
	if err != nil {
		return err
	}

	err = u.BCrypt.CompareAndHashPassword(user.Password, param.NewPassword)
	if err == nil {
		return errors.New("new password cannot be same as old password")
	}

	user.Password = hashPassword

	err = u.UserRepository.UpdateUser(tx, user)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
