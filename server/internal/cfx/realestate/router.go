package realestate

import (
	"degrens/panel/internal/routes"
	"degrens/panel/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RealEstateRouter struct {
	*routes.Router
}

func NewRealEstateRouter(rg *gin.RouterGroup) *RealEstateRouter {
	router := &RealEstateRouter{
		&routes.Router{
			RouterGroup: rg.Group("/realestate"),
		},
	}
	router.RegisterRoutes()
	return router
}

func (RR *RealEstateRouter) RegisterRoutes() {
	RR.RouterGroup.GET("/owned/:cid", RR.FetchOwnedRealEstate)
}

func (RR *RealEstateRouter) FetchOwnedRealEstate(ctx *gin.Context) {
	s_cid := ctx.Param("cid")

	cid, err := strconv.ParseInt(s_cid, 10, 32)
	log := logrus.WithField("cid", s_cid)
	if err != nil {
		log.WithError(err).Error("Failed to convert citizenid to uint")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the character you are trying to fetch",
		})
		return
	}

	locations, err := FetchLocationsForCitizenID(uint(cid))

	if err != nil {
		log.WithError(err).Error("Failed to fetch locations for citizen")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch the locations for the citizen",
		})
		return
	}

	ctx.JSON(200, locations)
}
