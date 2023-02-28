package reports

import (
	"degrens/panel/internal/auth/authinfo"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/users"
	"degrens/panel/lib/errors"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ReportRouter struct {
	routes.Router
}

type NewReportBody struct {
	Title string `json:"title"`
	// array of steamIds
	Members []string `json:"members"`
	// Array of tag names
	Tags []string `json:"tags"`
}

type FetchReportsBody struct {
	Offset int      `form:"offset"`
	Filter string   `form:"filter"`
	Tags   []string `form:"tags[]"`
	Open   bool     `form:"open"`
	Closed bool     `form:"closed"`
}

type UpdateReportTokensBody struct {
	ID uint `json:"id"`
	// ReportTag ids
	Tags []string `json:"tags"`
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
	RR.RouterGroup.POST("/new", RR.NewReportHandler())

	RR.RouterGroup.GET("/all", RR.FetchReports)

	RR.RouterGroup.GET("/:id", RR.FetchReportHandler())

	RR.RouterGroup.POST("/:id/member/:steamid", RR.HandleNewReportMember())

	RR.RouterGroup.GET("/tags", RR.FetchTagHandler())
	RR.RouterGroup.POST("/tags", RR.ReportTagHandler())
	RR.RouterGroup.PUT("/tags", RR.NewTagHandler())
}

func (RR *ReportRouter) FetchReports(ctx *gin.Context) {
	var body FetchReportsBody
	if err := ctx.ShouldBind(&body); err != nil {
		RR.Logger.Error("Failed to read body on GET /staff/reports request", "error", err)
		ctx.JSON(500, errors.BodyParsingFailed)
		return
	}

	userInfoPtr, exists := ctx.Get("userInfo")
	if !exists {
		ctx.JSON(403, errors.Unauthorized)
		return
	}
	userInfo := userInfoPtr.(*authinfo.AuthInfo)

	reportList, err := FetchReports(body.Filter, body.Offset, body.Tags, body.Open, body.Closed, userInfo)
	if err != nil {
		RR.Logger.Error("Failed to retrieve reports", "error", err, "filter", body)

		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Database Error",
			Description: "Seems like we had an error while fetching the reports in our database",
		})
		return
	}
	var reportTotal int64
	reportTotal, err = FetchReportCount(body.Filter, body.Offset, body.Tags, body.Open, body.Closed)
	if err != nil {
		RR.Logger.Error("Failed to retrieve report count", "error", err, "filter", body)

		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Database Error",
			Description: "Seems like we had an error while fetching the report count in our database",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"reports": reportList,
		"total":   reportTotal,
	})
}

func (RR *ReportRouter) NewReportHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body NewReportBody
		if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
			RR.Logger.Error("Failed to read body on a POST request to /staff/reports/new", "error", err)
			c.JSON(500, errors.BodyParsingFailed)
			return
		}
		userInfo, err := users.GetUserInfo(c)
		if err != nil {
			RR.Logger.Error("Failed to retrieve user info", "error", err)
			c.JSON(404, errors.Unauthorized)
			return
		}
		token, err := CreateNewReport(userInfo.Username, body.Title, body.Members, body.Tags, RR.Logger)
		if err != nil {
			RR.Logger.Error("Failed to create new report", "error", err)
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
			RR.Logger.Error("Failed to convert reportId to uint", "error", err, "id", ctx.Param("id"))
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to identify the report you are trying to fetch",
			})
			return
		}

		report, err := FetchReport(uint(reportId))
		if err != nil {
			RR.Logger.Error("Failed to retrieve reports", "error", err, "id", reportId)

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
			RR.Logger.Error("Failed to convert reportId to uint", "error", err, "id", ctx.Param("id"))
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
			RR.Logger.Error("Failed to get user info", "error", err, "id", ctx.Param("id"))
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to identify the report you are adding a member to",
			})
			return
		}
		userId, err := users.GetUserIdentifierForCtx(ctx)
		if err != nil {
			RR.Logger.Error("Failed to parse user info to a identifier", "error", err, "id", ctx.Param("id"))
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to identify the report you are adding a member to",
			})
			return
		}
		err = AddMemberToReport(userId, uint(reportId), steamId)
		if err != nil {
			RR.Logger.Error("Failed to add member to report", "error", err.Error(), "steamId", steamId, "reportId", reportId)
		}
		ctx.JSON(200, gin.H{})
	}
}

func (RR *ReportRouter) FetchTagHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tags, err := GetTags()
		if err != nil {
			RR.Logger.Error("Failed to retrieve tags", "error", err)
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Server Error",
				Description: "We encountered an error while trying to retrieve the report tags",
			})
			return
		}
		ctx.JSON(200, &tags)
	}
}

func (RR *ReportRouter) ReportTagHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body UpdateReportTokensBody
		if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil {
			RR.Logger.Error("Failed to read body on a POST request to /staff/reports/tags", "error", err)
			c.JSON(500, errors.BodyParsingFailed)
			return
		}
		if err := UpdateReportTags(body.ID, body.Tags); err != nil {
			RR.Logger.Error("Failed to update tags for report", "error", err)
			c.JSON(500, models.RouteErrorMessage{
				Title:       "Server Error",
				Description: "We encountered an error while updating the tags for the report",
			})
			return
		}
		c.JSON(200, gin.H{})
	}
}

func (RR *ReportRouter) NewTagHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body panel_models.ReportTag
		if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
			RR.Logger.Error("Failed to read body on a PUT request to /staff/reports/tags", "error", err)
			ctx.JSON(500, errors.BodyParsingFailed)
			return
		}
		if err := NewReportTag(body.Name, body.Color); err != nil {
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Server error",
				Description: "We encountered an error while adding a new tag",
			})
		}
		ctx.JSON(200, gin.H{})
	}
}
