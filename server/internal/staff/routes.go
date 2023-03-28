package staff

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/middlewares/role"
	"degrens/panel/internal/cfx"
	"degrens/panel/internal/cfx/business"
	"degrens/panel/internal/cfx/penalties"
	"degrens/panel/internal/cfx/players"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/staff/reportmessages"
	"degrens/panel/internal/staff/reports"
	"degrens/panel/internal/users"
	"degrens/panel/lib/errors"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StaffRouter struct {
	routes.Router
}

func NewStaffRouter(rg *gin.RouterGroup, logger log.Logger) {
	router := &StaffRouter{
		routes.Router{
			RouterGroup: rg.Group("/staff", role.New([]string{"staff"})),
			Logger:      logger,
		},
	}
	reportRG := rg.Group("/staff", role.New([]string{"staff", "player"}))
	reports.NewReportRouter(reportRG, logger)
	reportmessages.NewReportRouter(reportRG, logger)
	InitStaffService(logger)
	router.RegisterRoutes()
}

func (SR *StaffRouter) RegisterRoutes() {
	SR.RouterGroup.GET("/dashboard", SR.DashboardHandler())
	SR.RouterGroup.GET("/info/players", SR.FetchCfxPlayersHandler())

	SR.RouterGroup.GET("/notes", SR.FetchStaffNotes)
	SR.RouterGroup.POST("/notes", SR.createStaffNote)
	SR.RouterGroup.POST("/notes/:id", SR.updateStaffNote)
	SR.RouterGroup.DELETE("/notes/:id", SR.deleteStaffNote)
	SR.RouterGroup.GET("/logs", SR.fetchCxLogs)

	SR.RouterGroup.GET("/ban/list", SR.FetchBanList)
	SR.RouterGroup.POST("/ban/:id", SR.updateBan)
	SR.RouterGroup.DELETE("/ban/:id", SR.removeBan)

	business.NewBusinessRouter(SR.RouterGroup, SR.Logger)
	players.NewPlayerRouter(SR.RouterGroup, SR.Logger)
}

func (SR *StaffRouter) DashboardHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		info, err := GetDashboardInfo()
		if err != nil {
			if err.Code == 0 {
				err.Code = 500
			}
			c.JSON(err.Code, err.Message)
		} else {
			c.JSON(200, info)
		}
	}
}

func (SR *StaffRouter) FetchCfxPlayersHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// If the request takes long than 20s it is automatically yeeted, we force reset the mutex lock to prevent poising
		resetTimer := time.AfterFunc(20*time.Second, cfx.UnlockCache)
		plys, err := cfx.GetCfxPlayers()
		resetTimer.Stop()
		if err != nil {
			SR.Logger.Error("Failed to fetch information about cfx players", "error", err.Error())
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Server Error",
				Description: "We encountered an error while trying to fetch information about the players on the fiveM server",
			})
			return
		}
		ctx.JSON(200, plys)
	}
}

func (SR *StaffRouter) FetchStaffNotes(ctx *gin.Context) {
	notes, err := GetAllNotes()
	if err != nil {
		SR.Logger.Error("Failed to retrieve staff notes", "error", err)

		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server Error",
			Description: "Seems like we had an error while fetching the notes",
		})
		return
	}
	ctx.JSON(200, notes)
}

func (SR *StaffRouter) createStaffNote(ctx *gin.Context) {
	body := struct {
		Note string `json:"note"`
	}{}
	err := ctx.BindJSON(&body)
	if err != nil {
		SR.Logger.Error("Failed to parse staff create note body", "error", err)

		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "We couldn't create the staff note because we couldn't get the right data",
		})
		return
	}
	userInfoPtr, exists := ctx.Get("userInfo")
	if !exists {
		ctx.JSON(403, errors.Unauthorized)
		return
	}
	userInfo := userInfoPtr.(*authinfo.AuthInfo)
	fmt.Printf("%+v", body)
	err = CreateNote(userInfo.ID, body.Note)
	if err != nil {
		SR.Logger.Error("Failed to create staff note", "error", err)

		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server Error",
			Description: "Failed to create the staff note",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}

func (SR *StaffRouter) updateStaffNote(ctx *gin.Context) {
	body := struct {
		Note string `json:"note"`
	}{}
	err := ctx.BindJSON(&body)
	if err != nil {
		SR.Logger.Error("Failed to parse staff update note body", "error", err)

		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "We couldn't update the staff note because we couldn't get the right data",
		})
		return
	}
	noteIdStr := ctx.Param("id")
	if noteIdStr == "" {
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "We couldn't update the staff note because we couldn't determine the note id",
		})
		return
	}
	noteId, err := strconv.ParseUint(noteIdStr, 10, 64)
	if err != nil {
		SR.Logger.Error("Failed to parse note id from update note req", "error", err, "noteIdParam", noteIdStr)

		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "We couldn't update the staff note because we couldn't determine the note id",
		})
		return
	}

	userInfoPtr, exists := ctx.Get("userInfo")
	if !exists {
		ctx.JSON(403, errors.Unauthorized)
		return
	}
	userInfo := userInfoPtr.(*authinfo.AuthInfo)
	err = UpdateNote(userInfo.ID, uint(noteId), body.Note)
	if err != nil {
		SR.Logger.Error("Failed to update staff note", "error", err, "noteId", noteId)

		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server Error",
			Description: "Failed to update the staff note",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}

func (SR *StaffRouter) deleteStaffNote(ctx *gin.Context) {
	noteIdStr := ctx.Param("id")
	if noteIdStr == "" {
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "We couldn't update the staff note because we couldn't determine the note id",
		})
		return
	}
	noteId, err := strconv.ParseUint(noteIdStr, 10, 64)
	if err != nil {
		SR.Logger.Error("Failed to parse note id from update note req", "error", err, "noteIdParam", noteIdStr)

		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "We couldn't update the staff note because we couldn't determine the note id",
		})
		return
	}
	userInfoPtr, exists := ctx.Get("userInfo")
	if !exists {
		ctx.JSON(403, errors.Unauthorized)
		return
	}
	userInfo := userInfoPtr.(*authinfo.AuthInfo)
	err = DeleteNote(userInfo.ID, uint(noteId))
	if err != nil {
		SR.Logger.Error("Failed to delete staff note", "error", err, "noteId", noteId)

		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server Error",
			Description: "Failed to delete the staff note",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}

func (SR *StaffRouter) FetchBanList(ctx *gin.Context) {
	list, err := penalties.GetBanList()
	if err != nil {
		SR.Logger.Error("Failed to fetch the cfx ban list", "error", err)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server Error",
			Description: "Failed to fetch the banlist",
		})
		return
	}
	ctx.JSON(200, list)
}

func (SR *StaffRouter) updateBan(ctx *gin.Context) {
	banIdStr := ctx.Param("id")
	if banIdStr == "" {
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "No ban id given to update",
		})
	}
	banId, err := strconv.ParseUint(banIdStr, 10, 64)
	if err != nil {
		SR.Logger.Error("Failed to convert the banId string to a uint", "error", err, "banId", banIdStr)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "Could not transform your banId to a valid number",
		})
		return
	}

	body := struct {
		Reason string `json:"reason"`
		Length int    `json:"length"`
		Points uint   `json:"points"`
	}{}
	err = ctx.ShouldBindJSON(&body)
	if err != nil {
		SR.Logger.Error("Failed to bind the body of an update ban request", "error", err, "banId", banIdStr)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "Could get all the needed info from your request",
		})
		return
	}

	userId, err := users.GetUserIdentifierForCtx(ctx)
	if err != nil {
		SR.Logger.Error("Failed to get a user identifier string", "error", err, "banId", banIdStr)
		ctx.JSON(403, models.RouteErrorMessage{
			Title:       "Authentication Error",
			Description: "Could get valid identification for your request",
		})
		return
	}

	err = penalties.UpdateBan(userId, uint(banId), body.Points, body.Length, body.Reason)
	if err != nil {
		SR.Logger.Error("Failed to delete a cfx ban", "error", err)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server Error",
			Description: "Failed to delete the requested ban",
		})
		return
	}

	ctx.JSON(200, gin.H{})
}

func (SR *StaffRouter) removeBan(ctx *gin.Context) {
	banIdStr := ctx.Param("id")
	if banIdStr == "" {
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "No ban id given to delete",
		})
	}

	banId, err := strconv.ParseUint(banIdStr, 10, 64)
	if err != nil {
		SR.Logger.Error("Failed to convert the banId string to a uint", "error", err, "banId", banIdStr)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request Error",
			Description: "Could not transform your banId to a valid number",
		})
		return
	}

	userId, err := users.GetUserIdentifierForCtx(ctx)
	if err != nil {
		SR.Logger.Error("Failed to get a user identifier string", "error", err, "banId", banIdStr)
		ctx.JSON(403, models.RouteErrorMessage{
			Title:       "Authentication Error",
			Description: "Could get valid identification for your request",
		})
		return
	}

	err = penalties.RemoveBan(userId, uint(banId))
	if err != nil {
		SR.Logger.Error("Failed to delete a cfx ban", "error", err)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server Error",
			Description: "Failed to delete the requested ban",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}

func (SR *StaffRouter) fetchCxLogs(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "0")
	query := ctx.DefaultQuery("query", "*")
	page := uint64(0)
	if pageStr != "" {
		var err error
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			SR.Logger.Error("Failed to convert page to uint", "error", err, "page", pageStr)
			ctx.JSON(500, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to parse the current page number",
			})
			return
		}
	}

	logs, total, err := FetchCfxLogs(int(page), query)
	if err != nil {
		SR.Logger.Error("Failed to fetch panel logs from graylog", "error", err, "page", page)
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
