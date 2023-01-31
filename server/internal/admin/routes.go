package admin

import (
	"degrens/panel/internal/auth/middlewares/role"
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DevRouter struct {
	routes.Router
}

func NewDevRouter(rg *gin.RouterGroup, logger *log.Logger) {
	router := &DevRouter{
		routes.Router{
			RouterGroup: rg.Group("/dev", role.New("developer")),
			Logger:      *logger,
		},
	}
	router.RegisterRoutes()
}

func (DV *DevRouter) RegisterRoutes() {
	DV.RouterGroup.GET("/logs", DV.fetchPanelLogs)
}

func (DV *DevRouter) fetchPanelLogs(ctx *gin.Context) {
	offsetStr, ok := ctx.Params.Get("offset")
	offset := uint64(0)
	if ok {
		var err error
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			DV.Logger.Error("Failed to convert fetch offset to uint", "error", err, "offset", ctx.Param("offset"))
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to parse the fetch offset",
			})
			return
		}
	}
	logs, err := FetchPanelLogs(int(offset), "*")
	if err != nil {
		DV.Logger.Error("Failed to fetch panel logs from graylog", "error", err, "offset", offset)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch the panel logs",
		})
		return
	}
	ctx.JSON(200, logs)
}
