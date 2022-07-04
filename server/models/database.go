package models

import (
	"api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&User{}, &UserData{})
	if err != nil {
		return err
	}

	return nil
}
