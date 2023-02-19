package state

import (
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"

	"github.com/gin-gonic/gin"
)

type StateRouter struct {
	routes.Router
}

func NewStateRouter(rg *gin.RouterGroup, logger log.Logger) {
	router := &StateRouter{
		routes.Router{
			RouterGroup: rg.Group("/state"),
			Logger:      logger,
		},
	}
	router.RegisterRoutes()
}

func (ST *StateRouter) RegisterRoutes() {
	ST.RouterGroup.GET("/schedule/update", ST.FetchScheduledUpdate)
	ST.RouterGroup.POST("/schedule/update", ST.ScheduleUpdate)
}

func (ST *StateRouter) ScheduleUpdate(ctx *gin.Context) {
	scheduleUpdate()
	ctx.JSON(200, gin.H{})
}

func (ST *StateRouter) FetchScheduledUpdate(ctx *gin.Context) {
	ctx.JSON(200, incomingUpdate)
}
