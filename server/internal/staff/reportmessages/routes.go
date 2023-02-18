package reportmessages

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/staff/reports"
	"degrens/panel/lib/errors"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"net/http"
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
	RR.RouterGroup.POST("/message/add", RR.AddMessage)
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

func (RR *ReportRouter) AddMessage(ctx *gin.Context) {
	body := struct {
		Message  interface{} `json:"message"`
		ReportId uint        `json:"reportId"`
	}{}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		RR.Logger.Error("Failed to get the bind the body to the designated struct", "error", err)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to get the information from the request",
		})
		return
	}
	clientInfoPtr, exists := ctx.Get("userInfo")
	clientInfo := clientInfoPtr.(*authinfo.AuthInfo)
	if !exists {
		RR.Logger.Error("Failed to retrieve userinfo when joining report room")
		ctx.JSON(http.StatusForbidden, errors.Unauthorized)
		return
	}
	reportMsg, err := saveMessage(body.ReportId, body.Message, clientInfo)
	if err != nil {
		RR.Logger.Error("Failed to save a new report message", "error", err, "message", body.Message, "reportId", body.ReportId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "Failed to save the new report message",
		})
		return
	}
	ctx.JSON(200, reportMsg)
}
