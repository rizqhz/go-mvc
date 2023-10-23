package models

import (
	"github.com/rizghz/api/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(config *configs.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.MySqlConnectStr()), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, db.Error
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Blog{},
		&Book{},
	)
}
