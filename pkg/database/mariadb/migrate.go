package mariadb

import (
	"itfest-2025/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.Role{},
		&entity.User{},
		&entity.OtpCode{},
		&entity.Team{},
		&entity.Competition{},
		&entity.Announcement{},
		&entity.TeamProgress{},
		&entity.Submission{},
		&entity.TeamMember{},
		&entity.Registration{},
	)
	if err != nil {
		return err
	}

	return nil
}
