package reportmessages

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/staff/reports"
	"degrens/panel/lib/errors"
	"degrens/panel/models"
	"net/http"
	"strconv"

	"github.com/aidenwallis/go-utils/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ReportRouter struct {
	routes.Router
}

func NewReportRouter(rg *gin.RouterGroup) {
	router := &ReportRouter{
		routes.Router{
			RouterGroup: rg.Group("/reports"),
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
		logrus.WithField("id", ctx.Param("id")).WithError(err).Error("Failed to convert reportId to uint")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the report you are trying to join",
		})
		return
	}
	reportData, err := reports.FetchReport(uint(reportId))
	if err != nil {
		logrus.WithField("reportId", reportId).WithError(err).Error("Failed to convert fetch report by reportid")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to fetch the report you are trying to join",
		})
		return
	}
	report := reports.CreateReport(reportData)
	room := GetRoom(report)
	JoinReportRoom(ctx, room)
}

func (RR *ReportRouter) AddMessage(ctx *gin.Context) {
	body := struct {
		Message  map[string]interface{} `json:"message"`
		ReportId uint                   `json:"reportId"`
	}{}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		logrus.WithError(err).Error("Failed to get the bind the body to the designated struct")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to get the information from the request",
		})
		return
	}
	clientInfoPtr, exists := ctx.Get("userInfo")
	clientInfo := clientInfoPtr.(*authinfo.AuthInfo)
	if !exists {
		logrus.Error("Failed to retrieve userinfo when joining report room")
		ctx.JSON(http.StatusForbidden, errors.Unauthorized)
		return
	}

	report, err := reports.GetReport(body.ReportId)
	log := logrus.WithFields(logrus.Fields{
		"message":  body.Message,
		"reportId": body.ReportId,
	})
	if err != nil {
		log.WithError(err).Error("Failed to get report data")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "Failed to save the new report message, failed struct",
		})
		return
	}

	if body.Message == nil || len(utils.MapKeys(body.Message)) == 0 {
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Bad request",
			Description: "You cannot send an empty message",
		})
		return
	}
	reportMsg, err := report.AddMessage(body.ReportId, body.Message, clientInfo)
	if err != nil {
		log.WithError(err).Error("Failed to save a new report message")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "Failed to save the new report message",
		})
		return
	}
	err = SeedReportMessageMember(reportMsg)
	if err != nil {
		log.WithError(err).Error("Failed to seed report message sender")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "Failed to seed the new report message",
		})
		return
	}
	ctx.JSON(200, reportMsg)
}
