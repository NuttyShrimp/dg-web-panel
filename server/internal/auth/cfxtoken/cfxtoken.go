package cfxtoken

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/authinfo"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var tokenIds map[uint]string
var tokenInfoStore map[string]*TokenInfo

type TokenInfo struct {
	SteamId   string   `json:"steamId"`
	DiscordId string   `json:"discordId"`
	Username  string   `json:"username"`
	Roles     []string `json:"roles"`
}

func AuthorizeToken(ctx *gin.Context, token string) error {
	id := len(tokenIds)
	tokenIds[uint(id)] = token
	tokenInfo := TokenInfo{}
	ai, err := api.CfxApi.DoRequest("GET", fmt.Sprintf("/token/info/%s", token), nil, &tokenInfo)
	if err != nil {
		return err
	}
	if ai.Message != "" {
		return errors.New(fmt.Sprintf("Token error: %s", ai.Message))
	}
	tokenInfoStore[token] = &tokenInfo

	authInfo := authinfo.AuthInfo{
		ID:         uint(id),
		Roles:      tokenInfo.Roles,
		AuthMethod: "cfxtoken",
	}
	cookieSet := authInfo.Assign(ctx)
	if cookieSet != nil {
		return cookieSet
	}
	return nil
}

func RemoveToken(id uint) error {
	token, exists := tokenIds[id]
	if !exists {
		return errors.New("No token assigned to user Id")
	}
	delete(tokenIds, id)
	delete(tokenInfoStore, token)
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
	return info
}
