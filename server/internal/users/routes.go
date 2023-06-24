package users

import (
	"degrens/panel/internal/routes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserRouter struct {
	routes.Router
}

func NewUserRouter(rg *gin.RouterGroup) {
	router := &UserRouter{
		routes.Router{
			RouterGroup: rg.Group("/user"),
		},
	}
	router.RegisterRoutes()
}

func (UR *UserRouter) RegisterRoutes() {
	UR.RouterGroup.GET("/me", UR.meHandler())
}

// Sends back: username, avatarUrl, roles (in string)
func (UR *UserRouter) meHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := GetUserInfo(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not retrieve user info",
			})
			logrus.WithError(err).Error("Failed to retrieve user info")
			return
		}
		c.JSON(http.StatusOK, userInfo)
	}
}
