package players

import (
	"degrens/panel/internal/auth/middlewares/role"
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"degrens/panel/models"

	"github.com/gin-gonic/gin"
)

type PlayerRouter struct {
	routes.Router
}

func NewPlayerRouter(rg *gin.RouterGroup, logger *log.Logger) {
	router := &PlayerRouter{
		routes.Router{
			RouterGroup: rg.Group("/player", role.New([]string{"staff"})),
			Logger:      *logger,
		},
	}

	router.RegisterRoutes()
}

func (PR *PlayerRouter) RegisterRoutes() {
	PR.RouterGroup.GET("/:steamId/penalties", PR.FetchPlayerBanned)
}

func (PR *PlayerRouter) FetchPlayerBanned(ctx *gin.Context) {
	steamId := ctx.Param("steamId")
	until, err := IsPlayerBanned(steamId)
	if err != nil {
		PR.Logger.Error("Failed to fetch player ban status", "error", err, "steamId", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "Failed to fetch player ban status",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"until": until,
	})
}
