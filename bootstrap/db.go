package bootstrap

import (
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/config"
	"goblog/pkg/model"
	"time"

	"gorm.io/gorm"
)

func SetupDB() {
	db := model.ConnectDB()

	sqlDB, _ := db.DB()

	sqlDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))

	sqlDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))

	sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds")) * time.Second)
	// fmt.Println(123123)
	migration(db)

}

func migration(db *gorm.DB) {
	db.AutoMigrate(
		&user.User{},
		&article.Article{},
	)
}
