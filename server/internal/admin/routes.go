package admin

import (
	"degrens/panel/internal/auth/middlewares/role"
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var logger log.Logger

type DevRouter struct {
	routes.Router
}

func NewDevRouter(rg *gin.RouterGroup, logger2 log.Logger) {
	logger = logger2
	router := &DevRouter{
		routes.Router{
			RouterGroup: rg.Group("/dev", role.New([]string{"developer"})),
			Logger:      logger2,
		},
	}
	router.RegisterRoutes()
}

func (DV *DevRouter) RegisterRoutes() {
	DV.RouterGroup.GET("/logs", DV.fetchPanelLogs)
}

func (DV *DevRouter) fetchPanelLogs(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	page := uint64(0)
	if pageStr != "" {
		var err error
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			DV.Logger.Error("Failed to convert page to uint", "error", err, "page", pageStr)
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to parse the current page number",
			})
			return
		}
	}

	logs, total, err := FetchPanelLogs(int(page), "*")
	if err != nil {
		DV.Logger.Error("Failed to fetch panel logs from graylog", "error", err, "page", page)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch the panel logs",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"logs":  logs,
		"total": total,
	})
}
