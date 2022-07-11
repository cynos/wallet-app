package persistence

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfigPostgreSQL struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	ApplicationName string
	ConnectTimeout  int
	MaxOpenConn     int
	MaxIdleConn     int
}

func (config *DBConfigPostgreSQL) dbinfo() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s connect_timeout=%d application_name=%s",
		config.Host, config.Port, config.User, config.Name, config.Password, config.ConnectTimeout, config.ApplicationName,
	)
}

func (config *DBConfigPostgreSQL) Connect() (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(postgres.Open(config.dbinfo()), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	genericdb, _ := db.DB()
	return db, genericdb
}
