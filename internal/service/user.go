package service

import (
	"errors"
	"itfest-2025/entity"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/bcrypt"
	"itfest-2025/pkg/jwt"
	"itfest-2025/pkg/supabase"
	"mime/multipart"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(param *model.UserRegister) (string, error)
	UploadPayment(userID string, file *multipart.FileHeader) (string, error)
}

type UserService struct {
	UserRepository repository.IUserRepository
	TeamRepository repository.ITeamRepository
	BCrypt         bcrypt.Interface
	JwtAuth        jwt.Interface
	Supabase       supabase.Interface
}

func NewUserService(userRepository repository.IUserRepository, teamRepository repository.ITeamRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface) IUserService {
	return &UserService{
		UserRepository: userRepository,
		TeamRepository: teamRepository,
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
		PaymentTransc: "",
		RoleID:        2,
	}

	_, err = u.UserRepository.CreateUser(user)
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

func (u *UserService) UploadPayment(userID string, file *multipart.FileHeader) (string, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return "", err
	}

	user, err := u.UserRepository.GetUserByID(parsedID)
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
