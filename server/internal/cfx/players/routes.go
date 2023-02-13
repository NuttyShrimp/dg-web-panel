package players

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/middlewares/role"
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"degrens/panel/lib/utils"
	"degrens/panel/models"

	"github.com/gin-gonic/gin"
)

type PlayerRouter struct {
	routes.Router
}

func NewPlayerRouter(rg *gin.RouterGroup, logger log.Logger) {
	router := &PlayerRouter{
		routes.Router{
			RouterGroup: rg.Group("/player", role.New([]string{"staff"})),
			Logger:      logger,
		},
	}

	router.RegisterRoutes()
}

func (PR *PlayerRouter) RegisterRoutes() {
	PR.RouterGroup.GET("/:steamId/banstatus", PR.FetchPlayerBanned)
	PR.RouterGroup.GET("/:steamId/active", PR.FetchActiveCid)

	PR.RouterGroup.GET("/:steamId/penalties", PR.FetchPenalties)
	PR.RouterGroup.POST("/:steamId/warn", PR.WarnPlayer)
	PR.RouterGroup.POST("/:steamId/kick", PR.KickPlayer)
	PR.RouterGroup.POST("/:steamId/ban", PR.BanPlayer)
}

func (PR *PlayerRouter) getSteamIdFromParam(ctx *gin.Context) (string, bool) {
	steamId := ctx.Param("steamId")
	if !utils.ValidateSteamId(steamId) {
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request error",
			Description: "The steamid does not conform to the following format: steamid:\\d{15}",
		})
		return "", false
	}
	return steamId, true
}

func (PR *PlayerRouter) FetchPlayerBanned(ctx *gin.Context) {
	steamId := ctx.Param("steamId")
	until, err := IsPlayerBanned(steamId)
	if err != nil {
		PR.Logger.Error("Failed to fetch player ban status", "error", err, "steamId", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "Failed to fetch player ban status",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"until": until,
	})
}

func (PR *PlayerRouter) FetchActiveCid(ctx *gin.Context) {
	steamId, ok := PR.getSteamIdFromParam(ctx)
	if !ok {
		return
	}
	actCid, err := GetActiveCharacter(steamId)
	if err != nil {
		PR.Logger.Error("Failed to fetch characters for steamid", "error", err, "steamid", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to fetch the user characters",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"cid": actCid,
	})
}

func (PR *PlayerRouter) FetchPenalties(ctx *gin.Context) {
	steamId, ok := PR.getSteamIdFromParam(ctx)
	if !ok {
		return
	}
	list, err := GetPlayerPenalties(steamId)
	if err != nil {
		PR.Logger.Error("Failed to fetch penalties for steamid", "error", err, "steamid", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to fetch the user penalties",
		})
		return
	}
	ctx.JSON(200, list)
}

func (PR *PlayerRouter) WarnPlayer(ctx *gin.Context) {
	steamId, ok := PR.getSteamIdFromParam(ctx)
	if !ok {
		return
	}
	body := WarnInfo{}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		PR.Logger.Error("Failed to bind body for warn action", "error", err, "target", steamId)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to read the requests body",
		})
		return
	}
	ai, err := api.CfxApi.DoRequest("POST", "admin/actions/warn", &body, nil)
	if err != nil {
		PR.Logger.Error("Failed to warn the player on fivem server", "error", err, "target", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to warn the player from the server",
		})
		return
	}
	if ai.Message != "" {
		PR.Logger.Error("Failed to warn the player on fivem server", "error", ai.Message, "target", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to warn the player from the server",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}

func (PR *PlayerRouter) KickPlayer(ctx *gin.Context) {
	steamId, ok := PR.getSteamIdFromParam(ctx)
	if !ok {
		return
	}
	body := KickInfo{}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		PR.Logger.Error("Failed to bind body for kick action", "error", err, "target", steamId)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to read the requests body",
		})
		return
	}
	ai, err := api.CfxApi.DoRequest("POST", "admin/actions/kick", &body, nil)
	if err != nil {
		PR.Logger.Error("Failed to kick player on fivem server", "error", err, "target", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to kick the player from the server",
		})
		return
	}
	if ai.Message != "" {
		PR.Logger.Error("Failed to kick player on fivem server", "error", ai.Message, "target", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to kick the player from the server",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}

func (PR *PlayerRouter) BanPlayer(ctx *gin.Context) {
	steamId, ok := PR.getSteamIdFromParam(ctx)
	if !ok {
		return
	}
	body := BanInfo{}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		PR.Logger.Error("Failed to bind body for ban action", "error", err, "target", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to read the requests body",
		})
		return
	}
	ai, err := api.CfxApi.DoRequest("POST", "admin/actions/ban", &body, nil)
	if err != nil {
		PR.Logger.Error("Failed to ban player on fivem server", "error", err, "target", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to ban the player from the server",
		})
		return
	}
	if ai.Message != "" {
		PR.Logger.Error("Failed to ban player on fivem server", "error", ai.Message, "target", steamId)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to ban the player from the server",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}
