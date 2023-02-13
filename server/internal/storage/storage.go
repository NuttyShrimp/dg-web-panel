package storage

import (
	"degrens/panel/internal/config"
	"degrens/panel/lib/log"
)

type Storage interface {
	Add(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Remove(key string) error
	Move(key string, newKey string) error
	Clear()
	String() string
}

func InitStorages(conf *config.Config, logger log.Logger) {
	InitStateTokenStorage()
	InitCookieStore(conf, logger)
}
