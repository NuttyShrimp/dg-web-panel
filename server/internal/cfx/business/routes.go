package business

import (
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var br BusinessRouter

type BusinessRouter struct {
	routes.Router
}

func NewBusinessRouter(rg *gin.RouterGroup, logger *log.Logger) {
	br = BusinessRouter{
		Router: routes.Router{
			RouterGroup: rg.Group("/business"),
			Logger:      *logger,
		},
	}
	br.RegisterRoutes()
}

func (BR *BusinessRouter) RegisterRoutes() {
	BR.RouterGroup.GET("/all", BR.FetchAll)
	BR.RouterGroup.GET("/:id/logcount", BR.FetchLogCount)
	BR.RouterGroup.GET("/:id/logs", BR.FetchLogs)
	BR.RouterGroup.GET("/:id/employees", BR.FetchEmployees)
}

func (BR *BusinessRouter) FetchAll(ctx *gin.Context) {
	cidParam := ctx.Param("cid")
	filter := Filters{}
	if cidParam != "" {
		cid64, err := strconv.ParseUint(cidParam, 10, 32)
		if err != nil {
			BR.Logger.Error("Failed to convert cid to uint", "error", err, "cid", cid64)
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to convert the given CID to a valid number",
			})
			return
		}
		cid := uint(cid64)
		filter.cid = &cid
	}
	list, err := FetchBusinesses(&filter)
	if err != nil {
		BR.Logger.Error("Failed to fetch all businesses", "error", err)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch all businesses",
		})
		return
	}
	ctx.JSON(200, list)
}

func (BR *BusinessRouter) FetchLogCount(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		BR.Logger.Error("Failed to convert id to uint", "error", err, "id", ctx.Param("id"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	total, err := FetchLogCount(uint(id))
	if err != nil {
		BR.Logger.Error("Failed to fetch total amount of logs for a business", "error", err, "business_id", id)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch the total amount of logs for the selected business",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"total": total,
	})
}

func (BR *BusinessRouter) FetchLogs(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		BR.Logger.Error("Failed to convert id to uint", "error", err, "id", id)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	page := ctx.GetInt("page")
	if page < 0 {
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request error",
			Description: "page parameter cannot be negative",
		})
		return
	}
	logs, err := FetchLogs(uint(id), page)
	if err != nil {
		BR.Logger.Error("Failed to fetch logs for business", "error", err, "id", id, "page", page)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch the business logs",
		})
		return
	}
	ctx.JSON(200, logs)
}

func (BR *BusinessRouter) FetchEmployees(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		BR.Logger.Error("Failed to convert id to uint", "error", err, "id", id)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	list, err := FetchEmployees(uint(id))
	if err != nil {
		BR.Logger.Error("Failed to fetch employees for business", "error", err, "id", id)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch all employees for the selected business",
		})
		return
	}
	ctx.JSON(200, list)
}
