package flyers

import (
	"degrens/panel/internal/cfx/inventory"
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"degrens/panel/internal/routes"
	"degrens/panel/internal/users"
	panel_errors "degrens/panel/lib/errors"
	"degrens/panel/lib/graylogger"
	"degrens/panel/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FlyerRouter struct {
	routes.Router
}

func NewFlyerRouter(rg *gin.RouterGroup) {
	router := &FlyerRouter{
		routes.Router{
			RouterGroup: rg.Group("/flyers"),
		},
	}
	router.RegisterRoutes()
}

func (FR *FlyerRouter) RegisterRoutes() {
	FR.RouterGroup.GET("/", FR.FetchAllFlyers)
	FR.RouterGroup.GET("/retrieved", FR.FetchAllExistingFlyers)
	FR.RouterGroup.POST("/:id", FR.ApproveFlyer)
	FR.RouterGroup.DELETE("/:id", FR.RemoveFlyer)
}

func (FR *FlyerRouter) FetchAllFlyers(ctx *gin.Context) {
	flyers := []cfx_models.FlyerRequest{}
	err := db.CfxMariaDB.Client.Preload("Character.Info").Preload("Character.User").Find(&flyers).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Failed to fetch flyer requests")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while fetching the flyer requests",
		})
		return
	}
	ctx.JSON(200, flyers)
}

func (FR *FlyerRouter) FetchAllExistingFlyers(ctx *gin.Context) {
	metadata := map[string]string{"link": ""}
	items, err := inventory.FetchItemsByMetadata(metadata)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.WithError(err).Error("Failed to fetch flyer items")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while fetching the flyer items",
		})
		return
	}
	ctx.JSON(200, items)
}

func (FR *FlyerRouter) ApproveFlyer(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve a valid id from request path")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an issue while retrieving the flyer id you tried to approve",
		})
		return
	}

	userInfo, err := users.GetUserInfo(ctx)
	if err != nil {
		logrus.Error("Failed to get userinfo while giving vehicle to player")
		ctx.JSON(403, panel_errors.Unauthorized)
		return
	}

	flyerRequest := cfx_models.FlyerRequest{}
	db.CfxMariaDB.Client.Where("id = ?", id).First(&flyerRequest)
	flyerRequest.Approved = true
	db.CfxMariaDB.Client.Save(&flyerRequest)
	graylogger.Log("cfx:flyers:approve", fmt.Sprintf("%s(%d) has approved flyer request %d", userInfo.Username, userInfo.ID, id), "flyer", flyerRequest)

	ctx.JSON(200, gin.H{})
}

func (FR *FlyerRouter) RemoveFlyer(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve a valid id from request path")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an issue while retrieving the flyer id you tried to approve",
		})
		return
	}

	flyer := cfx_models.FlyerRequest{}
	err = db.CfxMariaDB.Client.First(&flyer, id).Error
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch flyer request")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, models.RouteErrorMessage{
				Title:       "Flyer not found",
				Description: "It seems like we couldn't find the flyer you tried to delete",
			})
			return
		}
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to get the flyer you tried to delete",
		})
	}
	db.CfxMariaDB.Client.Delete(&cfx_models.FlyerRequest{}, id)

	userInfo, err := users.GetUserInfo(ctx)
	if err != nil {
		logrus.Error("Failed to get userinfo while giving vehicle to player")
		ctx.JSON(403, panel_errors.Unauthorized)
		return
	}

	graylogger.Log("cfx:flyers:removed", fmt.Sprintf("%s(%d) has removed flyer request %d", userInfo.Username, userInfo.ID, id), "flyer", flyer)
}
