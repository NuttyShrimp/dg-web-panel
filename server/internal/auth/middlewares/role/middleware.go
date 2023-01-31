package role

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/users"
	"degrens/panel/lib/graylogger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxUserInfo, exists := ctx.Get("userInfo")
		if exists == false {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to retrieve userinfo",
			})
			return
		}
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)
		if !users.DoesUserHaveRole(userInfo.Roles, role) {
			graylogger.Log("auth:missing_role", fmt.Sprintf("user %d tried to access a link without having the right roles", userInfo.ID), "userInfo", userInfo, "requiredRole", role, "link", ctx.Request.URL)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
