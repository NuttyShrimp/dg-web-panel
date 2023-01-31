package users

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/db"
	"degrens/panel/internal/storage"
	"degrens/panel/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) (*models.UserInfo, error) {
	userInfoPtr, exists := c.Get("userInfo")
	if !exists {
		return nil, errors.New("Failed to retrieve user info from context")
	}
	userAuthInfo := userInfoPtr.(*authinfo.AuthInfo)
	userInfo := models.UserInfo{
		Roles: userAuthInfo.Roles,
	}
	if userAuthInfo.AuthMethod == authinfo.Discord {
		DBUser := db.MariaDB.Repository.GetUserById(userAuthInfo.ID)
		if &DBUser == nil {
			storage.RemoveCookie(c, "userInfo")
			return nil, errors.New("User not found")
		}
		userInfo.AvatarUrl = DBUser.AvatarUrl
		userInfo.Username = DBUser.Username
	} else {
	}
	return &userInfo, nil
}

func GetAuthInfo(userId uint) authinfo.AuthInfo {
	user := db.MariaDB.Repository.GetUserById(userId)
	registerdRoles := []string{}
	for _, role := range user.Roles {
		registerdRoles = append(registerdRoles, role.Name)
	}
	return authinfo.AuthInfo{
		ID:         userId,
		AuthMethod: authinfo.APIToken,
		Roles:      registerdRoles,
	}
}

func GetUserIdentifier(ctx *gin.Context) (string, error) {
	authInfoPtr, exists := ctx.Get("userInfo")
	if !exists {
		return "", errors.New("Failed to retrieve auth info for user from context")
	}
	authInfo := authInfoPtr.(*authinfo.AuthInfo)
	switch authInfo.AuthMethod {
	case "discord":
		{
			return strconv.Itoa(int(authInfo.ID)) + " (db)", nil
		}
	case "apitoken":
		{
			return strconv.Itoa(int(authInfo.ID)) + " (token)", nil
		}
	case "cfxtoken":
		{
			// TODO: get steamid from CFX API
			return strconv.Itoa(int(authInfo.ID)) + " (cfx)", nil
		}
	}
	return "", errors.New(fmt.Sprintf("failed to interpret authinfo for a valid user identifier: %+v", authInfo))
}
