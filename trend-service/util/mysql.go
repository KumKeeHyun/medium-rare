package util

import (
	"fmt"

	"github.com/KumKeeHyun/medium-rare/trend-service/config"
	"github.com/KumKeeHyun/medium-rare/trend-service/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func BuildMysqlConnection() (*gorm.DB, error) {
	urlFmt := "%s:%s@tcp(%s:%s)/trendDB?charset=utf8mb4&parseTime=True&loc=Local"
	url := fmt.Sprintf(urlFmt, config.App.MysqlConfig.User, config.App.MysqlConfig.Password, config.App.MysqlConfig.Host, config.App.MysqlConfig.Port)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.ReadRecord{})

	return db, nil
}
