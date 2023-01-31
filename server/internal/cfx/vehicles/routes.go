package vehicles

import (
	cfx_models "degrens/panel/internal/db/models/cfx"
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VehiclesRouter struct {
	routes.Router
}

func NewVehicleRouter(rg *gin.RouterGroup, logger *log.Logger) {
	vr := VehiclesRouter{
		Router: routes.Router{
			RouterGroup: rg.Group("/vehicles"),
			Logger:      *logger,
		},
	}
	vr.RegisterRoutes()
}

func (VR *VehiclesRouter) RegisterRoutes() {
	VR.RouterGroup.GET("/cid/:cid", VR.FetchAccounts)
}

func (VR *VehiclesRouter) FetchAccounts(ctx *gin.Context) {
	cid, err := strconv.ParseInt(ctx.Param("cid"), 10, 32)
	if err != nil {
		VR.Logger.Error("Failed to convert citizenid to uint", "error", err, "cid", ctx.Param("cid"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the character you are trying to fetch the vehicles for",
		})
		return
	}
	var vehs *[]cfx_models.PlayerVehicles
	vehs, err = FetchForCid(uint(cid))
	if err != nil {
		VR.Logger.Error("Failed to fetch vehicles for cid", "error", err, "cid", ctx.Param("cid"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: fmt.Sprintf("We encountered an error while trying to fetch the vehicles associated with cid: %d", cid),
		})
	}
	ctx.JSON(200, vehs)
}
