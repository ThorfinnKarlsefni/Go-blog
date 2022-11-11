package model

import (
	"goblog/pkg/logger"

	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ContentDB() *gorm.DB {
	var err error

	config := mysql.New(mysql.Config{
		DSN: "root:123456@tcp(127.0.0.1:3306)/goblog?charset=utf8&parseTime=True&loc=Local",
	})

	DB, err = gorm.Open(config, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	logger.LogError(err)

	return DB
}
