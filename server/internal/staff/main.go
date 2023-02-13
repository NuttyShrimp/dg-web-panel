package staff

import (
	"degrens/panel/internal/cfx"
	"degrens/panel/lib/log"
)

var logger log.Logger

func InitStaffService(pLogger log.Logger) {
	logger = pLogger
	go func() {
		_, err := cfx.GetCfxPlayers()
		if err != nil {
			logger.Error("Failed to retrieve information about CFX Players", "error", err)
		}
	}()
}
