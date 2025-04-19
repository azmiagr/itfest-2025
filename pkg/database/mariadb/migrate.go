package mariadb

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}
