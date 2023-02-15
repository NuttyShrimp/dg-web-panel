package cfxauth

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/cfxtoken"
	"fmt"

	"github.com/gin-gonic/gin"
)

func AuthorizeToken(ctx *gin.Context, token string) error {
	id := cfxtoken.GetNewToken()
	tokenInfo := cfxtoken.TokenInfo{}
	ai, err := api.CfxApi.DoRequest("GET", fmt.Sprintf("/tokens/info/%s", token), nil, &tokenInfo)
	if err != nil {
		return err
	}
	if ai.Message != "" {
		return fmt.Errorf("Token error: %s", ai.Message)
	}
	cfxtoken.RegisterToken(id, token, &tokenInfo)

	authInfo := authinfo.AuthInfo{
		ID:         uint(id),
		Roles:      append(tokenInfo.Roles, "player"),
		AuthMethod: "cfxtoken",
	}
	cookieSet := authInfo.Assign(ctx)
	if cookieSet != nil {
		return cookieSet
	}
	return nil
}
