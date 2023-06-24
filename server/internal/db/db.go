package db

import (
	"degrens/panel/internal/config"
	"degrens/panel/internal/db/mariadb"
	"degrens/panel/internal/db/redis"
)

var MariaDB *mariadb.Client
var Redis *redis.Client
var CfxMariaDB *mariadb.Client

func InitDatabase(conf *config.Config) {
	Redis = redis.InitRedisClient(conf)
	MariaDB = mariadb.InitMariaDBClient(&conf.MariaDB.Panel)
	CfxMariaDB = mariadb.InitMariaDBClient(&conf.MariaDB.Cfx)
	MariaDB.MariaDBMigrate()
}
