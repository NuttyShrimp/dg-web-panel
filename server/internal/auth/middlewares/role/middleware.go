package role

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/users"
	"degrens/panel/lib/graylogger"
	"fmt"
	"net/http"

	"github.com/aidenwallis/go-utils/utils"
	"github.com/gin-gonic/gin"
)

func New(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxUserInfo, exists := ctx.Get("userInfo")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Failed to retrieve userinfo",
			})
			return
		}
		userInfo := ctxUserInfo.(*authinfo.AuthInfo)

		_, ok := utils.SliceFind(roles, func(role string) bool {
			return users.DoesUserHaveRole(userInfo.Roles, role)
		})
		if !ok {
			graylogger.Log("auth:missing_role", fmt.Sprintf("user %d tried to access a link without having the right roles", userInfo.ID), "userInfo", userInfo, "requiredRoles", roles, "link", ctx.Request.URL)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
