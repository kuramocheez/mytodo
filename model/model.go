package model

import (
	"fmt"
	"mytodo/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitModel(config config.ProgramConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("Model: Tidak Dapat Terkoneksi Database, ", err.Error())
	}
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Users{}, &Category{})
}
