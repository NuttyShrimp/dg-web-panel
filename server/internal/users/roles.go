package users

import (
	"degrens/panel/internal/config"
	"errors"

	"github.com/gin-gonic/gin"
)

var roles *[]config.ConfigRole

func getRoleIndex(name string) int {
	for i, role := range *roles {
		if role.Name == name {
			return i
		}
	}
	return -1
}

func InitUserRoles(conf *config.Config) {
	roles = &conf.Discord.Roles
}

func HasRoleAccess(ctx *gin.Context, target string) (bool, error) {
	userInfo, err := GetUserInfo(ctx)
	if err != nil {
		return false, err
	}
	targetIdx := getRoleIndex(target)
	if targetIdx == -1 {
		return false, errors.New(target + " is an invalid role")
	}
	for _, role := range userInfo.Roles {
		if role == target {
			return true, nil
		}
		// Higher index = higher rank
		if getRoleIndex(role) <= targetIdx {
			return true, nil
		}
	}
	return false, nil
}

func DoesUserHaveRole(roles []string, role string) bool {
	roleIdx := getRoleIndex(role)
	for _, role := range roles {
		if getRoleIndex(role) <= roleIdx {
			return true
		}
	}
	return false
}
