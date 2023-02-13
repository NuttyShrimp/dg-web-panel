package reportmessages

import (
	"degrens/panel/internal/routes"
	"degrens/panel/internal/staff/reports"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportRouter struct {
	routes.Router
}

func NewReportRouter(rg *gin.RouterGroup, logger log.Logger) {
	router := &ReportRouter{
		routes.Router{
			RouterGroup: rg.Group("/reports"),
			Logger:      logger,
		},
	}
	router.RegisterRoutes()
}

func (RR *ReportRouter) RegisterRoutes() {
	RR.RouterGroup.GET("/join/:id", RR.ReportWS)
}

func (RR *ReportRouter) ReportWS(ctx *gin.Context) {
	reportId, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		RR.Logger.Error("Failed to convert reportId to uint", "error", err, "id", ctx.Param("id"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the report you are trying to join",
		})
		return
	}
	reportData, err := reports.FetchReport(uint(reportId))
	if err != nil {
		RR.Logger.Error("Failed to convert fetch report by reportid", "error", err, "reportId", reportId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to fetch the report you are trying to join",
		})
		return
	}
	report := reports.CreateReport(reportData)
	room := GetRoom(&report, RR.Logger)
	JoinReportRoom(ctx, room)
}
