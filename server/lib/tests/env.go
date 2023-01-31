package tests

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/discord"
	"degrens/panel/internal/config"
	"degrens/panel/internal/storage"
	"degrens/panel/lib/log"
)

var loadedEnvs = map[string]bool{
	"bare":    false,
	"graylog": false,
}

type Env struct {
	Config *config.Config
	Logger log.Logger
}

func LoadBareEnv() *Env {
	if loadedEnvs["bare"] {
		return nil
	}

	testConfig, err := config.LoadConfig("../../config.test.yml")

	if err != nil {
		panic(err)
	}
	env := Env{
		Logger: log.New(true),
		Config: testConfig,
	}

	discord.InitDiscordConf(testConfig, &env.Logger)
	storage.InitStorages(testConfig, &env.Logger)

	loadedEnvs["bare"] = true
	return &env
}

func LoadGraylogEnv() *Env {
	env := LoadBareEnv()
	if loadedEnvs["graylog"] {
		return nil
	}

	api.CreateGraylogApi(&env.Config.Graylog, &env.Logger)

	loadedEnvs["graylog"] = true
	return env
}
