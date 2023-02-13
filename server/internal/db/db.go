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

func InitDatabase(conf *config.Config, logger log.Logger) {
	Redis = redis.InitRedisClient(conf, logger)
	MariaDB = mariadb.InitMariaDBClient(&conf.MariaDB.Panel, logger)
	CfxMariaDB = mariadb.InitMariaDBClient(&conf.MariaDB.Cfx, logger)
	MariaDB.MariaDBMigrate()
}
