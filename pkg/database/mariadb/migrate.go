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
		&entity.Competition{},
		&entity.Team{},
		&entity.Announcement{},
		&entity.TeamProgress{},
		&entity.Submission{},
		&entity.TeamMember{},
	)
	if err != nil {
		return err
	}

	return nil
}
