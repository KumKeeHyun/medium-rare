package util

import (
	"github.com/KumKeeHyun/medium-rare/user-service/config"
	"github.com/KumKeeHyun/medium-rare/user-service/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func BuildMysqlConnection() (*gorm.DB, error) {
	// temp url
	config.App.MysqlConfig.DbURL = "root:balns@tcp(220.70.2.5:8081)/userDB?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(config.App.MysqlConfig.DbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.User{})

	return db, nil
}
