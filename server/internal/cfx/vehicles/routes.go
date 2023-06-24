package vehicles

import (
	"degrens/panel/internal/api"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/users"
	"degrens/panel/lib/errors"
	"degrens/panel/lib/graylogger"
	"degrens/panel/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type VehiclesRouter struct {
	routes.Router
}

func NewVehicleRouter(rg *gin.RouterGroup) {
	vr := VehiclesRouter{
		Router: routes.Router{
			RouterGroup: rg.Group("/vehicles"),
		},
	}
	vr.RegisterRoutes()
}

func (VR *VehiclesRouter) RegisterRoutes() {
	VR.RouterGroup.GET("/cid/:cid", VR.FetchAccounts)
	VR.RouterGroup.POST("/give", VR.GiveNewVehicle)
}

func (VR *VehiclesRouter) FetchAccounts(ctx *gin.Context) {
	cid, err := strconv.ParseInt(ctx.Param("cid"), 10, 32)
	if err != nil {
		logrus.WithField("cid", ctx.Param("cid")).WithError(err).Error("Failed to convert citizenid to uint")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the character you are trying to fetch the vehicles for",
		})
		return
	}
	var vehs *[]cfx_models.PlayerVehicles
	vehs, err = FetchForCid(uint(cid))
	if err != nil {
		logrus.WithField("cid", ctx.Param("cid")).WithError(err).Error("Failed to fetch vehicles for cid")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: fmt.Sprintf("We encountered an error while trying to fetch the vehicles associated with cid: %d", cid),
		})
	}
	ctx.JSON(200, vehs)
}

func (VR *VehiclesRouter) GiveNewVehicle(ctx *gin.Context) {
	body := struct {
		Model string `json:"model"`
		Owner uint   `json:"owner"`
	}{}
	err := ctx.BindJSON(&body)
	if err != nil {
		logrus.WithError(err).Error("Failed to bind body to struct when createing business")
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request error",
			Description: "We encountered an issue while reading the info from your request",
		})
		return
	}
	userInfo, err := users.GetUserInfo(ctx)
	if err != nil {
		logrus.Error("Failed to get userinfo while giving vehicle to player")
		ctx.JSON(403, errors.Unauthorized)
		return
	}
	ai, err := api.CfxApi.DoRequest("POST", "/vehicles/give", &body, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to give vehicle to player")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to give the vehicle to the player",
		})
		return
	}
	if ai.Message != "" {
		logrus.WithField("error", ai.Message).Error("Failed to give vehicle to player")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to give the vehicle to the player",
		})
		return
	}
	graylogger.Log("dev:actions:giveVehicle", fmt.Sprintf("%d (%s) heeft een voertuig (%s) aan %d gegeven", userInfo.ID, userInfo.Username, body.Model, body.Owner), "model", body.Model, "owner", body.Owner)
	ctx.JSON(200, gin.H{})
}
