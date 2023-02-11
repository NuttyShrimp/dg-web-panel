package players

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"errors"
	"time"

	"gorm.io/gorm"
)

type WarnInfo struct {
	Reason string `json:"reason"`
}

type KickInfo struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
}

type BanInfo struct {
	Points int    `json:"points"`
	Reason string `json:"reason"`
	Length int    `json:"length"`
}

func GetPlayerPenalties(steamId string) ([]cfx_models.Penalties, error) {
	penalties := []cfx_models.Penalties{}
	err := db.CfxMariaDB.Client.Where("steamId = ?", steamId).Find(&penalties).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return penalties, nil
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

func KickPlayer(steamId string, info *KickInfo) (bool, error) {
	input := struct {
		KickInfo
		Target string `json:"target"`
	}{
		KickInfo: *info,
		Target:   steamId,
	}
	output := struct {
		Result bool `json:"result"`
	}{
		Result: false,
	}

	ai, err := api.CfxApi.DoRequest("POST", "/admin/actions/kick", &input, &output)
	if err != nil {
		return false, err
	}
	if ai.Message != "" {
		return false, errors.New(ai.Message)
	}

	return output.Result, nil
}
