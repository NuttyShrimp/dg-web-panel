package cfxtoken

import (
	"degrens/panel/internal/api"
	"errors"
	"fmt"
)

var tokenIds map[uint]string = make(map[uint]string)
var tokenInfoStore map[string]*TokenInfo = make(map[string]*TokenInfo)

type TokenInfo struct {
	SteamId   string   `json:"steamId"`
	DiscordId string   `json:"discordId"`
	Username  string   `json:"username"`
	Roles     []string `json:"roles"`
}

func GetNewToken() uint {
	id := uint(len(tokenIds))
	tokenIds[id] = ""
	return id
}

func IsTokenValid(id uint) bool {
	_, ok := tokenIds[id]
	return ok
}

func RegisterToken(id uint, token string, info *TokenInfo) {
	tokenIds[id] = token
	tokenInfoStore[token] = info
}

func RemoveToken(id uint) error {
	token, exists := tokenIds[id]
	if !exists {
		return errors.New("No token assigned to user Id")
	}
	delete(tokenIds, id)
	delete(tokenInfoStore, token)
	if token == "" {
		return nil
	}
	ai, err := api.CfxApi.DoRequest("DELETE", fmt.Sprintf("/token/%s", token), nil, nil)
	if err != nil {
		return err
	}
	if ai.Message != "" {
		return errors.New(ai.Message)
	}
	return nil
}

func GetInfoForToken(id uint) *TokenInfo {
	token, exists := tokenIds[id]
	if !exists {
		return nil
	}
	info, exists := tokenInfoStore[token]
	if !exists {
		return nil
	}
	return info
}
