package auth

import (
	"degrens/panel/internal/auth/apikeys"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/cfxtoken"
	"degrens/panel/internal/storage"
	"degrens/panel/internal/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewMiddleWare() gin.HandlerFunc {
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
			userInfo = users.GetApiTokenAuthInfo(apikey.UserID)
		} else {
			// get sessionID
			userInfo, err = authinfo.GetUserInfo(c)
			if err != nil || (userInfo.AuthMethod == authinfo.CFXToken && !cfxtoken.IsTokenValid(userInfo.ID)) {
				storage.RemoveCookie(c, "userInfo")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Could not get user info",
				})
				return
			}
		}
		if userInfo.ID == 0 {
			storage.RemoveCookie(c, "userInfo")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Save userinfo in context as pointer
		c.Set("userInfo", &userInfo)
		c.Next()
	}
}
