package inventory

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/routes"
	"degrens/panel/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type InventoryRouter struct {
	routes.Router
}

func NewInventoryRouter(rg *gin.RouterGroup) {
	router := &InventoryRouter{
		routes.Router{
			RouterGroup: rg.Group("/inventory"),
		},
	}
	router.RegisterRoutes()
}

func (IR *InventoryRouter) RegisterRoutes() {
	IR.RouterGroup.DELETE("/:id", IR.deleteItem)
}

func (IR *InventoryRouter) deleteItem(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "URL error",
			Description: "We encountered an issue while retrieving the item id you tried to delete",
		})
		return
	}

	ai, err := api.CfxApi.DoRequest(http.MethodDelete, fmt.Sprintf("/inventory/item/%s", id), nil, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete item on cfx")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to delete the item",
		})
		return
	}
	if ai.Message != "" {
		logrus.WithField("error", ai.Message).Error("Failed to delete an item on cfx")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to delete the item",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}
