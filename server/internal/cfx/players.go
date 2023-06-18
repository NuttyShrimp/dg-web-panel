package cfx

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"degrens/panel/models"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type PlayersCache struct {
	Data      []cfx_models.User
	UpdatedAt time.Time
}

type ActivePlayerInfo struct {
	CitizenId uint `json:"cid"`
	ServerId  int  `json:"serverId"`
}

func GetCfxPlayers() (*[]cfx_models.User, error) {
	resetTimer := time.AfterFunc(20*time.Second, UnlockCache)
	cache := getCache()
	cache.Mutex.Lock()
	// Fetch updated players
	updatedPly := []cfx_models.User{}
	err := db.CfxMariaDB.Client.Preload("Points").Where("(last_updated BETWEEN ? AND ?) AND created_at < ?", cache.Players.UpdatedAt, time.Now(), cache.Players.UpdatedAt).Find(&updatedPly).Error
	if err != nil {
		return nil, err
	}

	// Fetch New players
	newPlys := []cfx_models.User{}
	err = db.CfxMariaDB.Client.Preload("Points").Where("created_at BETWEEN ? AND ?", cache.Players.UpdatedAt, time.Now()).Find(&newPlys).Error
	if err != nil {
		return nil, err
	}

	cache.Players.Data = append(cache.Players.Data, updatedPly...)
	cache.Players.Data = append(cache.Players.Data, newPlys...)
	cache.Players.UpdatedAt = time.Now()
	cache.Mutex.Unlock()
	resetTimer.Stop()
	return &cache.Players.Data, nil
}

func GetCfxPlayerInfo(steamId string) (*cfx_models.User, error) {
	cache := getCache()
	for i := range cache.Players.Data {
		if cache.Players.Data[i].SteamId == steamId {
			return &cache.Players.Data[i], nil
		}
	}
	cache.Mutex.Lock()
	info := cfx_models.User{}
	err := db.CfxMariaDB.Client.Find(&info, steamId).Error
	if err != nil {
		return nil, err
	}
	cache.Players.Data = append(cache.Players.Data, info)
	cache.Mutex.Unlock()
	return &info, nil
}

func GetCfxRolesForSteamId(steamId string) ([]string, error) {
	output := []string{}
	ei, err := api.CfxApi.DoRequest("GET", "/info/player/roles/"+steamId, nil, &output)
	if err != nil {
		return nil, err
	}
	if ei != nil && ei.Message != "" {
		return nil, errors.New(ei.Message)
	}
	return output, nil
}

func GetCfxUserInfo(steamId string) (*models.UserInfo, error) {
	info, err := GetCfxPlayerInfo(steamId)
	if err != nil {
		return nil, err
	}
	var roles []string
	roles, err = GetCfxRolesForSteamId(steamId)
	if err != nil {
		return nil, err
	}
	return &models.UserInfo{
		Username: info.Name,
		Roles:    roles,
	}, nil
}

// Gets the cfxPlayer based on the serverId
func GetCfxUserFromId(id uint) (*cfx_models.User, error) {
	output := struct {
		SteamId string `json:"steamId"`
	}{}
	ei, err := api.CfxApi.DoRequest(http.MethodGet, fmt.Sprintf("/info/serverId/%d", id), nil, &output)
	if err != nil {
		return nil, err
	}
	if ei != nil && ei.Message != "" {
		return nil, errors.New("Failed to retrieve players from Cfx Server, Error: " + ei.Message)
	}
	if output.SteamId != "" {
		return GetCfxPlayerInfo(output.SteamId)
	}
	return nil, errors.New("Invalid or inactive cfx serverId")
}

func GetSteamIdFromDiscordId(discordId string) string {
	playersRef, err := GetCfxPlayers()
	if err != nil {
		return ""
	}
	idToSearch := "discord:" + discordId
	players := *playersRef
	for i := range players {
		if players[i].DiscordId == idToSearch {
			return players[i].SteamId
		}
	}
	return ""
}

func GetActivePlayers() ([]ActivePlayerInfo, error) {
	info := []ActivePlayerInfo{}
	ai, err := api.CfxApi.DoRequest("GET", "/info/active", nil, &info)
	if err != nil {
		return info, err
	}
	if ai.Message != "" {
		return info, errors.New(ai.Message)
	}
	return info, nil
}
