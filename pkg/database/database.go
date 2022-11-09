package database

import (
	"database/sql"
	"goblog/pkg/logger"
	"time"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Initialize() {
	initDB()
	createTables()
}

func initDB() {
	var err error

	config := mysql.Config{
		User:                 "root",
		Passwd:               "123456",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}

	DB, err = sql.Open("mysql", config.FormatDSN())
	logger.LogError(err)

	DB.SetMaxOpenConns(100)

	DB.SetMaxIdleConns(25)

	DB.SetConnMaxLifetime(5 * time.Minute)

	err = DB.Ping()
	logger.LogError(err)
}

func createTables() {
	createArticlesSQL := `create table if not exists articles(
		id bigint(20) primary key auto_increment not null,
		title varchar(255) collate utf8mb4_unicode_ci not null,
		body longtext collate utf8mb4_unicode_ci
	);`

	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}
