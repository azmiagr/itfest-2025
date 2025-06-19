package service

import (
	"itfest-2025/entity"
	"itfest-2025/internal/repository"
	"itfest-2025/model"
	"itfest-2025/pkg/database/mariadb"
	"itfest-2025/pkg/mail"
	"strings"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAnnouncementService interface {
	SendAnnouncement(req model.RequestAnnouncement) error
	GetAnnouncement() ([]*model.ResponseAnnouncement, error)
}

type AnnouncementService struct {
	db                     *gorm.DB
	UserRepository         repository.IUserRepository
	TeamRepository         repository.ITeamRepository
	AnnouncementRepository repository.IAnnouncementRepository
}

func NewAnnouncementService(userRepository repository.IUserRepository, teamRepository repository.ITeamRepository, announcementRepository repository.IAnnouncementRepository) IAnnouncementService {
	return &AnnouncementService{
		db:                     mariadb.Connection,
		UserRepository:         userRepository,
		TeamRepository:         teamRepository,
		AnnouncementRepository: announcementRepository,
	}
}

func (a *AnnouncementService) GetAnnouncement() ([]*model.ResponseAnnouncement, error) {
	var response []*model.ResponseAnnouncement
	data, err := a.AnnouncementRepository.GetAnnouncement()
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		response = append(response, &model.ResponseAnnouncement{
			AnnouncementID: v.AnnouncementID.String(),
			Message: v.Description,
		}) 
	}

	return response, nil
}

func (a *AnnouncementService) SendAnnouncement(req model.RequestAnnouncement) error {
	users, err := a.UserRepository.GetAllUser()

	if err != nil {
		return err
	} else if len(users) == 0 {
		return model.ErrUserRecordNotFound
	}

	tx := a.db.Begin()
	defer tx.Rollback()

	err = a.AnnouncementRepository.CreateAnnouncement(tx, entity.Announcement{
		AnnouncementID: uuid.New(),
		Title: "Announcement",
		Description: req.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	mailBody := `
		<!DOCTYPE html>
		<html lang="id">

		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Pengumuman Terbaru IT FEST 2025</title>
			<style type="text/css">
				/* Reset CSS dasar untuk klien email */
				body, table, td, a {
					-webkit-text-size-adjust: 100%;
					-ms-text-size-adjust: 100%;
				}
				table, td {
					mso-table-lspace: 0pt;
					mso-table-rspace: 0pt;
				}
				img {
					-ms-interpolation-mode: bicubic;
					border: 0;
					height: auto;
					line-height: 100%;
					outline: none;
					text-decoration: none;
				}
				body {
					height: 100% !important;
					margin: 0 !important;
					padding: 0 !important;
					width: 100% !important;
				}

				/* Styling untuk teks yang konsisten */
				.paragraph {
					font-family: Arial, sans-serif; /* Fallback for Changa */
					font-size: 16px;
					line-height: 1.5;
					color: #d1d1d1;
					text-align: left;
				}
				.header-text {
					font-family: Arial, sans-serif; /* Fallback for Changa */
					font-size: 24px;
					font-weight: bold;
					color: #ffffff;
				}
				.footer-text {
					font-family: Arial, sans-serif; /* Fallback for Changa */
					font-size: 12px;
					line-height: 1.5;
					color: #a0a0a0;
				}

				/* Media Queries untuk Responsivitas (Opsional, tidak semua klien email mendukung) */
				@media screen and (max-width: 600px) {
					.full-width-image {
						width: 100% !important;
						max-width: 100% !important;
					}
					.header-text {
						font-size: 20px !important; /* Ukuran font lebih kecil di mobile */
					}
					.paragraph {
						font-size: 14px !important; /* Ukuran font lebih kecil di mobile */
					}
					.content-padding {
						padding: 10px !important;
					}
				}
			</style>
		</head>

		<body style="margin: 0; padding: 0; background-color: #030D35;">
			<div style="background: linear-gradient(to bottom, #030D35 0%, #19217C 100%);">
				<table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px; margin: 0 auto; background-color: transparent;">
					<tr>
						<td align="center" valign="top" style="padding: 40px 20px 20px 20px;" class="content-padding">
							<table border="0" cellpadding="0" cellspacing="0" width="100%">

								<tr>
									<td align="center" style="padding-bottom: 20px;">
										<img src="https://i.postimg.cc/9QHJbbGw/it-fest-2025.png" width="300" alt="IT FEST 2025 Logo" style="display: block; width: 300px; max-width: 100%; min-width: 100px; font-family: Arial, sans-serif; color: #ffffff;" class="full-width-image">
									</td>
								</tr>

								<tr>
									<td align="center" style="padding: 20px 0;">
										<img src="https://i.postimg.cc/m2gbD3cM/email-ITFEST.png" width="130" alt="Ikon Pengumuman" style="display: block; width: 130px;">
									</td>
								</tr>

								<tr>
									<td align="center" style="padding: 10px 0; font-family: Arial, sans-serif; font-size: 24px; font-weight: bold; color: #ffffff;" class="header-text">
										Pengumuman Terbaru ITFEST 2025
									</td>
								</tr>

								<tr>
									<td align="center" style="padding: 10px 20px; font-family: Arial, sans-serif; font-size: 16px; line-height: 1.5; color: #d1d1d1; text-align: left;" class="paragraph">
										Halo Para Peserta, berikut pengumuman penting dari Panitia IT FEST 2025:
										<br><br>
											$MESSAGE$
										<br><br>
										Untuk informasi lebih lengkap, silakan kunjungi laman Dashboard Anda.
									</td>
								</tr>

								<tr>
									<td align="center" style="padding: 20px 20px; font-family: Arial, sans-serif; font-size: 16px; line-height: 1.5; color: #d1d1d1;" class="paragraph">
										Terima kasih,<br>
										Tim Panitia IT FEST
									</td>
								</tr>

								<tr>
									<td align="center" style="padding: 0 20px 40px 20px; font-family: Arial, sans-serif; font-size: 12px; line-height: 1.5; color: #a0a0a0;" class="footer-text">
										Keluarga Besar Mahasiswa Sistem Informasi
										<br>
										Universitas Brawijaya
									</td>
								</tr>

							</table>
						</td>
					</tr>
				</table>
			</div>
		</body>
		</html>
	`
	mailBody = strings.Replace(mailBody, "$MESSAGE$", req.Message, 1)

	for _, v := range users {
		if v.RoleID == 2 && v.StatusAccount == "active" {
			err = mail.SendEmail(v.Email, "Pengumuman IT FEST 2025", mailBody)
		}
	}

	if err != nil {
		return err
	}
	return nil
}
