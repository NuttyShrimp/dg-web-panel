package state

import (
	"degrens/panel/internal/routes"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type StateRouter struct {
	routes.Router
}

func NewStateRouter(rg *gin.RouterGroup) {
	router := &StateRouter{
		routes.Router{
			RouterGroup: rg.Group("/state"),
		},
	}
	router.RegisterRoutes()
}

func (ST *StateRouter) RegisterRoutes() {
	ST.RouterGroup.GET("/schedule/update", ST.FetchScheduledUpdate)
	ST.RouterGroup.POST("/schedule/update", ST.ScheduleUpdate)
}

func (ST *StateRouter) ScheduleUpdate(ctx *gin.Context) {
	wasSchedule := strconv.FormatBool(incomingUpdate)
	scheduleUpdate()
	logrus.Info("Scheduled an update")
	ctx.String(200, wasSchedule)
}

func (ST *StateRouter) FetchScheduledUpdate(ctx *gin.Context) {
	ctx.JSON(200, incomingUpdate)
}
