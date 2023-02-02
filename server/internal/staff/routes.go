package staff

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/middlewares/role"
	"degrens/panel/internal/cfx"
	"degrens/panel/internal/cfx/business"
	"degrens/panel/internal/cfx/players"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/staff/reportmessages"
	"degrens/panel/internal/staff/reports"
	"degrens/panel/lib/errors"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StaffRouter struct {
	routes.Router
}

func NewStaffRouter(rg *gin.RouterGroup, logger *log.Logger) {
	router := &StaffRouter{
		routes.Router{
			RouterGroup: rg.Group("/staff", role.New([]string{"staff"})),
			Logger:      *logger,
		},
	}
	reportRG := rg.Group("/staff", role.New([]string{"staff", "player"}))
	reports.NewReportRouter(reportRG, logger)
	reportmessages.NewReportRouter(reportRG, logger)
	InitStaffService(logger)
	router.RegisterRoutes()
}

func (SR *StaffRouter) RegisterRoutes() {
	// TODO: Move to secured endpoint which validates user roles to include staff or higher
	SR.RouterGroup.GET("/dashboard", SR.DashboardHandler())
	SR.RouterGroup.GET("/info/players", SR.FetchCfxPlayersHandler())

	SR.RouterGroup.GET("/notes", SR.FetchStaffNotes)
	SR.RouterGroup.POST("/notes", SR.createStaffNote)
	SR.RouterGroup.POST("/notes/:id", SR.updateStaffNote)
	SR.RouterGroup.DELETE("/notes/:id", SR.deleteStaffNote)

	business.NewBusinessRouter(SR.RouterGroup, &SR.Logger)
	players.NewPlayerRouter(SR.RouterGroup, &SR.Logger)
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
		return
	}
}

func (SR *StaffRouter) FetchCfxPlayersHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result := make(chan interface{})
		statusCode := make(chan int)
		go func() {
			plys, err := cfx.GetCfxPlayers()
			if err != nil {
				SR.Logger.Error("Failed to fetch information about cfx players", "error", err.Error())
				statusCode <- 500
				result <- models.RouteErrorMessage{
					Title:       "Server Error",
					Description: "We encountered an error while trying to fetch information about the players on the fiveM server",
				}
			}
			statusCode <- 200
			result <- &plys
		}()
		ctx.JSON(<-statusCode, <-result)
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
