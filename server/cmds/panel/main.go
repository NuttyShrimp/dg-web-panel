package main

import (
	"degrens/panel/cmds/router"
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/discord"
	"degrens/panel/internal/config"
	"degrens/panel/internal/db"
	"degrens/panel/internal/storage"
	"degrens/panel/internal/users"
	graylog "degrens/panel/lib/graylogger"
	"degrens/panel/lib/logs/sentryhook"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/getsentry/sentry-go"
)

func main() {
	// Load config
	configPath := config.GetConfigPath()
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	// Setup Sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://52914b9b7d394dc8962680061a942328@sentry.nuttyshrimp.me/4",
		Environment:      conf.Server.Env,
		EnableTracing:    true,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
	}); err != nil {
		logrus.Errorf("Sentry initialization failed: %v\n", err)
	}

	// Create logger
	logrus.AddHook(sentryhook.Hook{})
	if conf.Server.Env == "development" {
		logrus.SetLevel(logrus.TraceLevel)
	}

	db.InitDatabase(conf)
	graylog.InitGrayLogger(&conf.Graylog)
	api.CreateGraylogApi(&conf.Graylog)
	api.CreateCfxApi(&conf.Cfx)

	// Create discord auth config
	discord.InitDiscordConf(conf)
	users.InitUserRoles(conf)
	storage.InitStorages(conf)

	defer sentry.Flush(2 * time.Second)

	r := router.SetupRouter(conf)
	err = r.Run(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
	if err != nil {
		logrus.Errorf("Could not start server: %s", err)
		return
	}
}
