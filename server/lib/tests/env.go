package tests

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/discord"
	"degrens/panel/internal/config"
	"degrens/panel/internal/storage"
)

var loadedEnvs = map[string]bool{
	"bare":    false,
	"graylog": false,
}

type Env struct {
	Config *config.Config
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
		Config: testConfig,
	}

	discord.InitDiscordConf(testConfig)
	storage.InitStorages(testConfig)

	loadedEnvs["bare"] = true
	return &env
}

func LoadGraylogEnv() *Env {
	env := LoadBareEnv()
	if loadedEnvs["graylog"] {
		return nil
	}

	api.CreateGraylogApi(&env.Config.Graylog)

	loadedEnvs["graylog"] = true
	return env
}
