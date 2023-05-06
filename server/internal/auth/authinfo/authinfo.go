package authinfo

import (
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/internal/storage"
	"errors"

	"github.com/gin-gonic/gin"
)

type AuthInfo struct {
	ID         uint       `json:"id"`
	Roles      []string   `json:"roles"`
	AuthMethod AuthMethod `json:"authMethod"`
}

type AuthMethod string

const (
	Discord  AuthMethod = "discord"
	APIToken AuthMethod = "apitoken"
	CFXToken AuthMethod = "cfxtoken"
)

// get/set userinfo to/from cookies, data is stored in redis tied to UUID that is set in the cookie
func (AI *AuthInfo) Assign(c *gin.Context) error {
	// Set cookie
	cookieSet := storage.AddHiddenCookie(c, "userInfo", AI)
	if !cookieSet {
		return errors.New("Failed to set cookie with userInfo")
	}
	return nil
}

func GetAuthInfoFromUser(user *panel_models.User) *AuthInfo {
	return &AuthInfo{
		ID:         user.ID,
		Roles:      user.GetRoleNames(),
		AuthMethod: Discord,
	}
}

func GetUserInfo(c *gin.Context) (AuthInfo, error) {
	var userInfo AuthInfo
	err := storage.GetHiddenCookie(c, "userInfo", &userInfo)
	if err != nil {
		return AuthInfo{}, err
	}
	return userInfo, err
}
