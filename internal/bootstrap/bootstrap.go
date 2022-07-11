package bootstrap

import (
	"log"
	"os"

	"gorm.io/gorm"

	"github.com/wallet-app/internal/infrastructure/cache"
	"github.com/wallet-app/internal/infrastructure/persistence"
	"github.com/wallet-app/internal/tools"
)

var db *gorm.DB
var cacher *cache.CacheRedis

func init() {
	// init database
	dbconf := persistence.DBConfigPostgreSQL{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	db, _ = dbconf.Connect()

	// init redis
	redisClient := cache.RedisClient(cache.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		DB:       tools.StringsToInt(os.Getenv("REDIS_DB")),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	cacher = cache.NewCacheRedis(redisClient)
}

func LaunchApp() {
	switch os.Getenv("SERVICE") {
	case "EMONEY_SERVICE":
		serviceEmoney()
	case "TOPUP_SERVICE":
		serviceTopUp()
	default:
		log.Fatal("invalid service")
	}
}
