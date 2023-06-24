package reports

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/users"
	"degrens/panel/lib/errors"
	"degrens/panel/models"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type ReportRouter struct {
	routes.Router
}

type NewReportBody struct {
	Title string `json:"title"`
	// array of steamIds
	Members []string `json:"members"`
}

type FetchReportsBody struct {
	Offset int    `form:"offset"`
	Filter string `form:"filter"`
	Open   bool   `form:"open"`
	Closed bool   `form:"closed"`
}

type UpdateReportTokensBody struct {
	ID uint `json:"id"`
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
	RR.RouterGroup.POST("/new", RR.NewReportHandler())

	RR.RouterGroup.GET("/all", RR.FetchReports)

	RR.RouterGroup.GET("/:id", RR.FetchReportHandler())

	RR.RouterGroup.POST("/:id/member/:steamid", RR.HandleNewReportMember())
}

func (RR *ReportRouter) FetchReports(ctx *gin.Context) {
	var body FetchReportsBody
	if err := ctx.ShouldBind(&body); err != nil {
		logrus.WithError(err).Error("Failed to read body on GET /staff/reports request")
		ctx.JSON(500, errors.BodyParsingFailed)
		return
	}

	userInfoPtr, exists := ctx.Get("userInfo")
	if !exists {
		ctx.JSON(403, errors.Unauthorized)
		return
	}
	userInfo := userInfoPtr.(*authinfo.AuthInfo)

	reportList, err := FetchReports(body.Filter, body.Offset, body.Open, body.Closed, userInfo)
	if err != nil {
		logrus.WithField("filter", body).WithError(err).Error("Failed to retrieve reports")

		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Database Error",
			Description: "Seems like we had an error while fetching the reports in our database",
		})
		return
	}
	var reportTotal int64
	reportTotal, err = FetchReportCount(body.Filter, body.Offset, body.Open, body.Closed)
	if err != nil {
		logrus.WithField("filter", body).WithError(err).Error("Failed to retrieve report count")

		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Database Error",
			Description: "Seems like we had an error while fetching the report count in our database",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"reports": reportList,
		"total":   math.Ceil(float64(reportTotal) / 25),
	})
}

func (RR *ReportRouter) NewReportHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body NewReportBody
		if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
			logrus.WithError(err).Error("Failed to read body on a POST request to /staff/reports/new")
			c.JSON(500, errors.BodyParsingFailed)
			return
		}
		userInfo, err := users.GetUserInfo(c)
		if err != nil {
			logrus.WithError(err).Error("Failed to retrieve user info")
			c.JSON(404, errors.Unauthorized)
			return
		}
		token, err := CreateNewReport(userInfo.Username, body.Title, body.Members)
		if err != nil {
			logrus.WithError(err).Error("Failed to create new report")
			c.JSON(500, models.RouteErrorMessage{
				Title:       "Server Error",
				Description: "We encountered an error while creating a new report",
			})
			return
		}
		c.JSON(200, gin.H{
			"token": token,
		})
	}
}

func (RR *ReportRouter) FetchReportHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reportId, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			logrus.WithField("id", ctx.Param("id")).WithError(err).Error("Failed to convert reportId to uint")
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to identify the report you are trying to fetch",
			})
			return
		}

		report, err := FetchReport(uint(reportId))
		if err != nil {
			logrus.WithField("id", reportId).WithError(err).Error("Failed to retrieve reports")

			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Database Error",
				Description: "Seems like we had an error while fetching the reports in our database",
			})
			return
		}
		ctx.JSON(200, gin.H{
			"report": report,
		})
	}
}

func (RR *ReportRouter) HandleNewReportMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reportId, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
		if err != nil {
			logrus.WithField("id", ctx.Param("id")).WithError(err).Error("Failed to convert reportId to uint")
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to identify the report you are adding a member to",
			})
			return
		}
		steamId := ctx.Param("steamid")
		if steamId == "" {
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "URL error",
				Description: "No valid steamId was given to add to the a new member to the report",
			})
			return
		}
		if err != nil {
			logrus.WithField("id", ctx.Param("id")).WithError(err).Error("Failed to get user info")
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to identify the report you are adding a member to",
			})
			return
		}
		userId, err := users.GetUserIdentifierForCtx(ctx)
		if err != nil {
			logrus.WithField("id", ctx.Param("id")).WithError(err).Error("Failed to parse user info to a identifier")
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to identify the report you are adding a member to",
			})
			return
		}
		err = AddMemberToReport(userId, uint(reportId), steamId)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"steamId":  steamId,
				"reportId": reportId,
			}).WithError(err).Error("Failed to add member to report")
		}
		ctx.JSON(200, gin.H{})
	}
}
