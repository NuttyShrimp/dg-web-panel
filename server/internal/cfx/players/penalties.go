package players

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"time"

	"gorm.io/gorm"
)

func GetPlayerPenalties(steamId string) ([]cfx_models.Penalties, error) {
	penalties := []cfx_models.Penalties{}
	err := db.CfxMariaDB.Client.Where("steamId = ?", steamId).Find(&penalties).Error
	return penalties, err
}

func IsPlayerBanned(steamId string) (*time.Time, error) {
	penalty := cfx_models.Penalties{}
	err := db.CfxMariaDB.Client.Where("steamId = ? AND penalty = ?", steamId, cfx_models.BanPenalty).First(&penalty).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if penalty.Length == -1 {
		startTime := time.Unix(0, 0)
		return &startTime, nil
	}
	until := penalty.Date.Add(time.Duration(penalty.Length) * time.Second)
	return &until, nil
}
