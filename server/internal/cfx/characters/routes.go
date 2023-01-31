package characters

import (
	"degrens/panel/internal/cfx/bank"
	"degrens/panel/internal/cfx/vehicles"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CharacterRouter struct {
	routes.Router
}

func NewCharacterRouter(rg *gin.RouterGroup, logger *log.Logger) {
	router := &CharacterRouter{
		routes.Router{
			RouterGroup: rg.Group("/character"),
			Logger:      *logger,
		},
	}
	router.RegisterRoutes()
}

func (CR *CharacterRouter) RegisterRoutes() {
	CR.RouterGroup.GET("/:cid", CR.ValidateCID)
	CR.RouterGroup.GET("/all", CR.FetchAllCharacters)
	CR.RouterGroup.GET("/:cid/info", CR.FetchCharInfo)
	CR.RouterGroup.GET("/:cid/reputation", CR.FetchCharRep)
	bank.NewBankRouter(CR.RouterGroup, &CR.Logger)
	vehicles.NewVehicleRouter(CR.RouterGroup, &CR.Logger)
}

func (CR *CharacterRouter) ValidateCID(ctx *gin.Context) {
	cid, err := strconv.ParseInt(ctx.Param("cid"), 10, 32)
	if err != nil {
		CR.Logger.Error("Failed to convert citizenid to uint", "error", err, "cid", ctx.Param("cid"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the character you are trying to fetch",
		})
		return
	}

	statusCode := 200
	if !DoesCIDExist(uint(cid)) {
		statusCode = 404
	}
	ctx.JSON(statusCode, gin.H{})
}

func (CR *CharacterRouter) FetchAllCharacters(ctx *gin.Context) {
	chars, err := GetAllCharacters()
	if err != nil {
		CR.Logger.Error("Failed to fetch all characters", "error", err, "id", ctx.Param("id"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while fetch all characters you are trying to fetch",
		})
		return
	}
	ctx.JSON(200, chars)
}

func (CR *CharacterRouter) FetchCharInfo(ctx *gin.Context) {
	cid, err := strconv.ParseInt(ctx.Param("cid"), 10, 32)
	if err != nil {
		CR.Logger.Error("Failed to convert citizenid to uint", "error", err, "cid", ctx.Param("cid"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the character you are trying to fetch",
		})
		return
	}
	char := cfx_models.Character{}
	char, err = GetCharacterInfo(uint(cid))
	if err != nil {
		CR.Logger.Error("Failed to fetch character info", "error", err, "cid", cid)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: fmt.Sprint("We encountered an error while trying to fetch the info for the character with cid %d", cid),
		})
		return
	}
	ctx.JSON(200, char)
}

func (CR *CharacterRouter) FetchCharRep(ctx *gin.Context) {
	cid, err := strconv.ParseInt(ctx.Param("cid"), 10, 32)
	if err != nil {
		CR.Logger.Error("Failed to convert citizenid to uint", "error", err, "cid", ctx.Param("cid"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the character you are trying to fetch reputation for",
		})
		return
	}
	rep := &cfx_models.CharacterReputation{}
	rep, err = GetCharacterReputation(uint(cid))

	if err != nil {
		CR.Logger.Error("Failed to fetch character reputation", "error", err, "cid", ctx.Param("cid"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to fetch the character reputation",
		})
		return
	}
	ctx.JSON(200, rep)
}
