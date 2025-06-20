package service

import (
	"errors"
	"fmt"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"
	"itfest-2025/pkg/mail"
	"time"

	"gorm.io/gorm"
)

type IOtpService interface {
	ResendOtp(param model.GetOtp) error
	ResendOtpChangePassword(param model.GetOtp) error
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
									<img src="https://i.postimg.cc/9QHJbbGw/it-fest-2025.png" width="300" alt="IT FEST 2025 Logo" style="display: block; width: 300px; max-width: 100%%; min-width: 100px; font-family: Arial, sans-serif; color: #ffffff;">
								</td>
							</tr>

							<tr>
								<td align="center" style="padding: 20px 0;">
									<img src="https://i.postimg.cc/pdCm4W3M/kode.png" width="200" alt="Email Icon" style="display: block; width: 200px;">
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
	`, otp.Code))
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

func (o *OtpService) ResendOtpChangePassword(param model.GetOtp) error {
	tx := o.db.Begin()
	defer tx.Rollback()

	user, err := o.UserRepository.GetUser(model.UserParam{
		UserID: param.UserID,
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

	if otp.UpdatedAt.After(time.Now().UTC().Add(-5 * time.Minute)) {
		return errors.New("you can only resend otp every 5 minutes")
	}

	otp.Code = mail.GenerateCode()

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
