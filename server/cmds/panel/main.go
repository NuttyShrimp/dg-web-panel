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
	"degrens/panel/lib/log"
	"fmt"
	"time"

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
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	// Create logger
	logger := log.New(conf.Server.Env == "development", "https://52914b9b7d394dc8962680061a942328@sentry.nuttyshrimp.me/4")

	db.InitDatabase(conf, logger)
	graylog.InitGrayLogger(&conf.Graylog, logger)
	api.CreateGraylogApi(&conf.Graylog, logger)
	api.CreateCfxApi(&conf.Cfx, logger)

	// Create discord auth config
	discord.InitDiscordConf(conf, logger)
	users.InitUserRoles(conf)
	storage.InitStorages(conf, logger)

	defer sentry.Flush(2 * time.Second)

	r := router.SetupRouter(conf, logger)
	err = r.Run(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
	if err != nil {
		logger.Errorf("Could not start server: %s", err)
		return
	}
}
