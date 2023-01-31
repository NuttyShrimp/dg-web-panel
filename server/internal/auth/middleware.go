package auth

import (
	"degrens/panel/internal/auth/apikeys"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/storage"
	"degrens/panel/internal/users"
	"degrens/panel/lib/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var logger log.Logger

func NewMiddleWare(logger2 *log.Logger) gin.HandlerFunc {
	logger = *logger2
	return func(c *gin.Context) {
		// Check for header this is highest in rank
		header := c.GetHeader("X-Api-Key")
		var err error
		var userInfo authinfo.AuthInfo
		if header != "" {
			apikey := apikeys.GetAPIKey(header)
			if apikey == nil || apikey.ApiKey == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Could not get user info",
				})
				return
			}
			userInfo = users.GetAuthInfo(apikey.UserID)
		} else {
			// get sessionID
			userInfo, err = authinfo.GetUserInfo(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Could not get user info",
				})
				return
			}
		}
		if &userInfo == nil || userInfo.ID == 0 {
			storage.RemoveCookie(c, "userInfo")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Save userinfo in context as pointer
		c.Set("userInfo", &userInfo)
		c.Next()
	}
}
