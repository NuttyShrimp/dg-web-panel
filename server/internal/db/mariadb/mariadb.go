package mariadb

import (
	"degrens/panel/internal/config"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/lib/log"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	Client     *gorm.DB
	logger     log.Logger
	Repository *Repository
}

func InitMariaDBClient(config *config.ConfigMariaConn, logger *log.Logger) *Client {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	retries := 0
	for err != nil {
		time.Sleep(5 * time.Second)
		if retries < 10 {
			(*logger).Errorf("Failed to connect to a mariadb instance(%s), try: %d, trying again...", config.Database, retries)
			time.Sleep(5 * time.Second)
			retries++
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			continue
		}
		(*logger).Fatalf("Failed to connect to a mariadb instance(%s) after %d tries: %s", retries, config.Database, err)
	}
	(*logger).Info("Connected to mariadb instance")
	return &Client{
		Client:     db,
		logger:     *logger,
		Repository: newRepository(db),
	}
}

func (m *Client) MariaDBMigrate() {
	err := m.Client.AutoMigrate(&panel_models.User{}, &panel_models.Role{}, &panel_models.DiscordTokens{}, &panel_models.Report{}, &panel_models.ReportMember{}, &panel_models.ReportMessage{}, &panel_models.ReportTag{}, &panel_models.APIKey{}, &panel_models.Notes{})
	if err != nil {
		m.logger.Fatalf("Failed to migrate database: %s", err)
		return
	}
}
