package cfx

import (
	"degrens/panel/internal/cfx/flyers"
	"degrens/panel/internal/cfx/inventory"
	"degrens/panel/internal/routes"

	"github.com/gin-gonic/gin"
)

type CfxRouter struct {
	routes.Router
}

func NewCfxRouter(rg *gin.RouterGroup) {
	router := &CfxRouter{
		routes.Router{
			RouterGroup: rg.Group("/cfx"),
		},
	}
	router.RegisterRoutes()
}

func (CR *CfxRouter) RegisterRoutes() {
	flyers.NewFlyerRouter(CR.RouterGroup)
	inventory.NewInventoryRouter(CR.RouterGroup)
}
