package util

import (
	"github.com/KumKeeHyun/medium-rare/user-service/config"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func BuildMysqlConnection() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.App.MysqlConfig.DbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.User{})

	return db, nil
}
