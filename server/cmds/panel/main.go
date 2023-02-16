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
	"os"
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
		Dsn:              "https://16ccf13e3a274fb9bcb6f827bd8f57d0@sentry.nuttyshrimp.me/11",
		Environment:      conf.Server.Env,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	defer sentry.Flush(2 * time.Second)

	// Create logger
	logger := log.New(conf.Server.Env == "development", "https://16ccf13e3a274fb9bcb6f827bd8f57d0@sentry.nuttyshrimp.me/11")

	db.InitDatabase(conf, logger)
	graylog.InitGrayLogger(&conf.Graylog, logger)
	api.CreateGraylogApi(&conf.Graylog, logger)
	api.CreateCfxApi(&conf.Cfx, logger)
	isValid := api.ValidateGraylogApi()
	if !isValid {
		os.Exit(1)
	}

	// Create discord auth config
	discord.InitDiscordConf(conf, logger)
	users.InitUserRoles(conf)
	storage.InitStorages(conf, logger)

	r := router.SetupRouter(conf, logger)
	err = r.Run(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port))
	if err != nil {
		logger.Errorf("Could not start server: %s", err)
		return
	}
}
