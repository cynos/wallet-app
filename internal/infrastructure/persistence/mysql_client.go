package persistence

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfigMySQL struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	TimeLocation string
	ParseTime    string
}

func (config *DBConfigMySQL) dbinfo() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=%s&loc=%s",
		config.User, config.Password, config.Host, config.Port, config.Name, config.ParseTime, config.TimeLocation,
	)
}

func (config *DBConfigMySQL) Connect() (*gorm.DB, *sql.DB) {
	fmt.Println(config.dbinfo())
	db, err := gorm.Open(mysql.Open(config.dbinfo()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	generic, _ := db.DB()
	return db, generic
}
