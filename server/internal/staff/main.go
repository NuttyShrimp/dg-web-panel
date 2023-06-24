package staff

import (
	"degrens/panel/internal/cfx"

	"github.com/sirupsen/logrus"
)

func InitStaffService() {
	go func() {
		_, err := cfx.GetCfxPlayers()
		if err != nil {
			logrus.WithError(err).Error("Failed to retrieve information about CFX Players")
		}
	}()
}
