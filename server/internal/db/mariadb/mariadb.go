package mariadb

import (
	"degrens/panel/internal/config"
	panel_models "degrens/panel/internal/db/models/panel"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	Client     *gorm.DB
	Repository *Repository
}

func InitMariaDBClient(conf *config.ConfigMariaConn) *Client {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	retries := 0
	for err != nil {
		time.Sleep(5 * time.Second)
		if retries < 10 {
			logrus.Errorf("Failed to connect to a mariadb instance(%s), try: %d, trying again...", conf.Database, retries)
			time.Sleep(5 * time.Second)
			retries++
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			continue
		}
		logrus.WithError(err).Fatalf("Failed to connect to a mariadb instance(%s) after %d tries", conf.Database, retries)
	}
	logrus.Info("Connected to mariadb instance")
	return &Client{
		Client:     db,
		Repository: newRepository(db),
	}
}

func (m *Client) MariaDBMigrate() {
	err := m.Client.AutoMigrate(&panel_models.User{}, &panel_models.Role{}, &panel_models.DiscordTokens{}, &panel_models.Report{}, &panel_models.ReportMember{}, &panel_models.ReportMessage{}, &panel_models.APIKey{}, &panel_models.Notes{})
	if err != nil {
		logrus.Fatalf("Failed to migrate database: %s", err)
		return
	}
}
