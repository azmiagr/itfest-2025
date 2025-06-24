package service

import (
	"errors"
	"fmt"
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
	UploadKTM(userID uuid.UUID, file *multipart.FileHeader) error
	VerifyUser(param model.VerifyUser) error
	UpdateProfile(userID uuid.UUID, param model.UpdateProfile) (*model.UpdateProfile, error)
	GetUserProfile(userID uuid.UUID) (model.UserProfile, error)
	GetMyTeamProfile(userID uuid.UUID) (*model.UserTeamProfile, error)
	ChangePassword(email string) (string, error)
	ChangePasswordAfterVerify(param model.ResetPasswordRequest) error
	VerifyOtpChangePassword(param model.VerifyToken) error
	CompetitionRegistration(userID uuid.UUID, competitionID int, param model.CompetitionRegistrationRequest) error
	GetUserPaymentStatus() ([]*model.GetUserPaymentStatus, error)
	GetTotalParticipant() (*model.GetTotalParticipant, error)
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

	if param.Password != param.ConfirmPassword {
		return result, errors.New("password doesn't match")
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

	token, err := u.JwtAuth.CreateJWTToken(user.UserID, false)
	if err != nil {
		return result, errors.New("failed to create token")
	}

	team := &entity.Team{
		TeamID:        uuid.New(),
		TeamName:      "",
		TeamStatus:    "belum terverifikasi",
		UserID:        user.UserID,
		CompetitionID: 1,
	}

	err = u.TeamRepository.CreateTeam(tx, team)
	if err != nil {
		return result, err
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

	err = mail.SendEmail(user.Email, "OTP Verification", fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<style>
				body, table, td, a {
					-webkit-text-size-adjust: 100%%;
					-ms-text-size-adjust: 100%%;
				}

				table, td {
					mso-table-lspace: 0pt;
					mso-table-rspace: 0pt;
				}

				img {
					-ms-interpolation-mode: bicubic;
					border: 0;
					height: auto;
					line-height: 100%%;
					outline: none;
					text-decoration: none;
				}

				body {
					height: 100%% !important;
					margin: 0 !important;
					padding: 0 !important;
					width: 100%% !important;
				}
			</style>
		</head>

		<body style="margin: 0; padding: 0; background-color: #030D35; background: linear-gradient(to bottom, #030D35 0%%, #19217C 100%%);">
			<table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 600px; margin: 0 auto;">
				<tr>
					<td align="center" valign="top" style="padding: 40px 20px 20px 20px;">
						<table border="0" cellpadding="0" cellspacing="0" width="100%%">

							<tr>
								<td align="center" style="padding-bottom: 20px;">
									<img src="https://i.imgur.com/3fcE9Ll.png" width="300" alt="IT FEST 2025 Logo" style="display: block; width: 300px; max-width: 100%%; min-width: 100px; font-family: Arial, sans-serif; color: #ffffff;">
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 20px 0;">
									<img src="https://i.imgur.com/dgvL3Gf.png" width="200" alt="Email Icon" style="display: block; width: 200px;">
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 10px 0; font-family: Arial, sans-serif; font-size: 24px; font-weight: bold; color: #ffffff;">
									Kode Verifikasi Anda
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 10px 20px; font-family: Arial, sans-serif; font-size: 16px; line-height: 1.5; color: #d1d1d1;">
									Gunakan kode di bawah ini untuk menyelesaikan proses verifikasi email Anda. Kode ini hanya berlaku selama 5 menit.
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 30px 0;">
									<table border="0" cellspacing="0" cellpadding="0" width="100%%" style="max-width: 576px;">
										<tr>
											<td align="center" style="border-radius: 8px; background-color: #072547; padding: 20px 25px;">
												<div style="font-family: Arial, sans-serif; font-size: 36px; font-weight: bold; color: #85FFF5; letter-spacing: 5px; text-shadow: 0px 0px 15px rgba(255,255,255,0.6);">
													%s
												</div>
											</td>
										</tr>
									</table>
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 30px 20px 20px 20px; font-family: Arial, sans-serif; font-size: 14px; line-height: 1.5; color: #a0a0a0;">
									Jika Anda tidak merasa mendaftar untuk IT FEST, abaikan saja email ini.
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 0 20px 40px 20px; font-family: Arial, sans-serif; font-size: 12px; line-height: 1.5; color: #a0a0a0 !important;">
									Keluarga Besar Mahasiswa Departemen Sistem Informasi<br>
									Universitas Brawijaya
								</td>
							</tr>

						</table>
					</td>
				</tr>
			</table>
		</body>
		</html>
		`, code))

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
	var isAdmin bool

	tx := u.db.Begin()
	defer tx.Rollback()

	var result model.LoginResponse

	user, err := u.UserRepository.GetUser(model.UserParam{
		Email: param.Email,
	})
	if err != nil {
		return result, errors.New("email or password is wrong")
	}

	if user.RoleID == 1 {
		isAdmin = true
	} else {
		isAdmin = false
	}

	err = u.BCrypt.CompareAndHashPassword(user.Password, param.Password)
	if err != nil {
		return result, errors.New("email or password is wrong")
	}

	token, err := u.JwtAuth.CreateJWTToken(user.UserID, isAdmin)
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
	maxSize := int64(1024 * 1024)
	if file.Size > maxSize {
		return "", errors.New("file size exceeds maximum limit of 1MB")
	}

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

func (u *UserService) UploadKTM(userID uuid.UUID, file *multipart.FileHeader) error {
	maxSize := int64(1024 * 1024)
	if file.Size > maxSize {
		return errors.New("file size exceeds maximum limit of 1MB")
	}

	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return err
	}

	ktmURL, err := u.Supabase.UploadFile(file)
	if err != nil {
		return err
	}

	user.StudentCardLink = ktmURL

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

func (u *UserService) UpdateProfile(userID uuid.UUID, param model.UpdateProfile) (*model.UpdateProfile, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: userID,
	})

	if err != nil {
		return nil, err
	}

	user.FullName = param.FullName
	user.StudentNumber = param.StudentNumber
	user.University = param.University
	user.Major = param.Major
	user.PhoneNumber = param.PhoneNumber

	err = u.UserRepository.UpdateUser(tx, user)
	if err != nil {
		return nil, err
	}

	response := &model.UpdateProfile{
		FullName:      user.FullName,
		StudentNumber: user.StudentNumber,
		University:    user.University,
		Major:         user.Major,
		PhoneNumber:   user.PhoneNumber,
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return response, nil
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
		TeamName:            team.TeamName,
		StudentNumber:       user.StudentNumber,
		Deadline:            competititon.Deadline,
		CompetitionCategory: competititon.CompetitionName,
		Members:             memberResponse,
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return TeamProfileResponse, nil

}

func (u *UserService) ChangePassword(email string) (string, error) {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.UserRepository.GetUser(model.UserParam{
		Email: email,
	})
	if err != nil {
		return "", err
	}

	otp := mail.GenerateCode()
	err = u.OtpRepository.CreateOtp(tx, &entity.OtpCode{
		OtpID:  uuid.New(),
		UserID: user.UserID,
		Code:   otp,
	})
	if err != nil {
		return "", err
	}

	err = mail.SendEmail(user.Email, "OTP Atur Ulang Kata Sandi", fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<style>
				body, table, td, a {
					-webkit-text-size-adjust: 100%%;
					-ms-text-size-adjust: 100%%;
				}

				table, td {
					mso-table-lspace: 0pt;
					mso-table-rspace: 0pt;
				}

				img {
					-ms-interpolation-mode: bicubic;
					border: 0;
					height: auto;
					line-height: 100%%;
					outline: none;
					text-decoration: none;
				}

				body {
					height: 100%% !important;
					margin: 0 !important;
					padding: 0 !important;
					width: 100%% !important;
				}
			</style>
		</head>

		<body style="margin: 0; padding: 0; background-color: #030D35; background: linear-gradient(to bottom, #030D35 0%%, #19217C 100%%);">
			<table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 600px; margin: 0 auto;">
				<tr>
					<td align="center" valign="top" style="padding: 40px 20px 20px 20px;">
						<table border="0" cellpadding="0" cellspacing="0" width="100%%">

							<tr>
								<td align="center" style="padding-bottom: 20px;">
									<img src="/public/images/it-fest-2025.png" width="300" alt="IT FEST 2025 Logo" style="display: block; width: 300px; max-width: 100%%; min-width: 100px; font-family: Arial, sans-serif; color: #ffffff;">
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 20px 0;">
									<img src="https://i.imgur.com/dgvL3Gf.png" width="200" alt="Email Icon" style="display: block; width: 200px;">
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 10px 0; font-family: Arial, sans-serif; font-size: 24px; font-weight: bold; color: #ffffff;">
									Kode Atur Ulang Kata Sandi Anda
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 10px 20px; font-family: Arial, sans-serif; font-size: 16px; line-height: 1.5; color: #d1d1d1;">
									Kami menerima permintaan untuk mengatur ulang kata sandi akun IT FEST Anda. Gunakan kode di bawah ini pada halaman yang tersedia. Kode ini hanya berlaku selama 5 menit.
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 30px 0;">
									<table border="0" cellspacing="0" cellpadding="0" width="100%%" style="max-width: 576px;">
										<tr>
											<td align="center" style="border-radius: 8px; background-color: #072547; padding: 20px 25px;">
												<div style="font-family: Arial, sans-serif; font-size: 36px; font-weight: bold; color: #85FFF5; letter-spacing: 5px; text-shadow: 0px 0px 15px rgba(255,255,255,0.6);">
													%s
												</div>
											</td>
										</tr>
									</table>
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 30px 20px 20px 20px; font-family: Arial, sans-serif; font-size: 14px; line-height: 1.5; color: #a0a0a0;">
									Jika Anda tidak merasa mendaftar untuk IT FEST, abaikan saja email ini.
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 0 20px 40px 20px; font-family: Arial, sans-serif; font-size: 12px; line-height: 1.5; color: #a0a0a0 !important;">
									Keluarga Besar Mahasiswa Departemen Sistem Informasi<br>
									Universitas Brawijaya
								</td>
							</tr>

						</table>
					</td>
				</tr>
			</table>
		</body>
		</html>
	`, otp))
	if err != nil {
		return "", err
	}

	jwtToken, err := u.JwtAuth.CreateJWTToken(user.UserID, false)
	if err != nil {
		return "", nil
	}

	err = tx.Commit().Error
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (u *UserService) VerifyOtpChangePassword(param model.VerifyToken) error {
	tx := u.db.Begin()
	defer tx.Rollback()

	otp, err := u.OtpRepository.GetOtp(tx, model.GetOtp{
		UserID: param.UserID,
		Code:   param.OTP,
	})
	if err != nil {
		return err
	}

	if otp.Code != param.OTP {
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

func (u *UserService) ChangePasswordAfterVerify(param model.ResetPasswordRequest) error {
	tx := u.db.Begin()
	defer tx.Rollback()

	user, err := u.UserRepository.GetUser(model.UserParam{
		UserID: param.UserID,
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

func (u *UserService) CompetitionRegistration(userID uuid.UUID, competitionID int, param model.CompetitionRegistrationRequest) error {
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
	user.PhoneNumber = param.PhoneNumber

	err = u.UserRepository.UpdateUser(tx, user)
	if err != nil {
		return err
	}

	team, err := u.TeamRepository.GetTeamByUserID(tx, userID)
	if err != nil {
		return err
	}

	team.CompetitionID = competitionID
	err = u.TeamRepository.UpdateTeam(tx, team)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUserPaymentStatus() ([]*model.GetUserPaymentStatus, error) {
	var res []*model.GetUserPaymentStatus

	tx := u.db.Begin()
	defer tx.Rollback()

	users, err := u.UserRepository.GetAllUser()
	if err != nil {
		return nil, err
	}

	for _, v := range users {
		competition, err := u.CompetitionRepository.GetCompetitionByID(tx, v.Team.CompetitionID)
		if err != nil {
			continue
		}
		res = append(res, &model.GetUserPaymentStatus{
			FullName:        v.FullName,
			StudentNumber:   v.StudentNumber,
			Email:           v.Email,
			PaymentTransc:   v.PaymentTransc,
			TeamName:        v.Team.TeamName,
			TeamStatus:      v.Team.TeamStatus,
			CompetitionName: competition.CompetitionName,
		})
	}

	return res, nil
}

func (u *UserService) GetTotalParticipant() (*model.GetTotalParticipant, error) {

	var (
		totalUIUX int
		totalBP   int
	)

	tx := u.db.Begin()
	defer tx.Rollback()

	users, err := u.UserRepository.GetAllUser()
	if err != nil {
		return nil, err
	}

	for _, v := range users {
		if v.Team.CompetitionID == 2 {
			totalUIUX++
		} else if v.Team.CompetitionID == 3 {
			totalBP++
		}
	}

	res := &model.GetTotalParticipant{
		TotalUIUX: totalUIUX,
		TotalBP:   totalBP,
	}

	return res, nil
}
