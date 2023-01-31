package users

import (
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	routes.Router
}

func NewUserRouter(rg *gin.RouterGroup, logger *log.Logger) {
	router := &UserRouter{
		routes.Router{
			RouterGroup: rg.Group("/user"),
			Logger:      *logger,
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
			UR.Logger.Error("Failed to retrieve user info", "error", err)
			return
		}
		c.JSON(http.StatusOK, userInfo)
	}
}
