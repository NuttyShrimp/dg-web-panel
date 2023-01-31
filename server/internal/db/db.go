package db

import (
	"degrens/panel/internal/config"
	"degrens/panel/internal/db/mariadb"
	"degrens/panel/internal/db/redis"
	"degrens/panel/lib/log"
)

var MariaDB *mariadb.Client
var Redis *redis.Client
var CfxMariaDB *mariadb.Client

func InitDatabase(config *config.Config, logger *log.Logger) {
	Redis = redis.InitRedisClient(config, logger)
	MariaDB = mariadb.InitMariaDBClient(&config.MariaDB.Panel, logger)
	CfxMariaDB = mariadb.InitMariaDBClient(&config.MariaDB.Cfx, logger)
	MariaDB.MariaDBMigrate()
}
