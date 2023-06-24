package business

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/authinfo"
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

var br BusinessRouter

type BusinessRouter struct {
	routes.Router
}

func NewBusinessRouter(rg *gin.RouterGroup) {
	br = BusinessRouter{
		Router: routes.Router{
			RouterGroup: rg.Group("/business"),
		},
	}
	br.RegisterRoutes()
}

func (BR *BusinessRouter) RegisterRoutes() {
	BR.RouterGroup.GET("/all", BR.FetchAll)
	BR.RouterGroup.GET("/:id/logcount", BR.FetchLogCount)
	BR.RouterGroup.GET("/:id/logs", BR.FetchLogs)
	BR.RouterGroup.GET("/:id/employees", BR.FetchEmployees)
	BR.RouterGroup.POST("/:id/owner", BR.ChangeOwner)

	// TODO: secure with role authenitction
	BR.RouterGroup.POST("/new", BR.CreateBusiness)
}

func (BR *BusinessRouter) FetchAll(ctx *gin.Context) {
	cidParam := ctx.Param("cid")
	filter := Filters{}
	if cidParam != "" {
		cid64, err := strconv.ParseUint(cidParam, 10, 32)
		if err != nil {
			logrus.WithField("cid", cid64).WithError(err).Error("Failed to convert cid to uint")
			ctx.JSON(400, models.RouteErrorMessage{
				Title:       "Parsing error",
				Description: "We encountered an error while trying to convert the given CID to a valid number",
			})
			return
		}
		cid := uint(cid64)
		filter.cid = &cid
	}
	list, err := FetchBusinesses(&filter)
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch all businesses")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch all businesses",
		})
		return
	}
	ctx.JSON(200, list)
}

func (BR *BusinessRouter) FetchLogCount(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithField("id", ctx.Param("id")).WithError(err).Error("Failed to convert id to uint")
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	total, err := FetchLogCount(uint(id))
	if err != nil {
		logrus.WithField("business_id", id).WithError(err).Error("Failed to fetch total amount of logs for a business")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch the total amount of logs for the selected business",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"total": total,
	})
}

func (BR *BusinessRouter) FetchLogs(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithField("id", id).WithError(err).Error("Failed to convert id to uint")
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	page, err := strconv.ParseUint(ctx.DefaultQuery("page", "0"), 10, 32)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":   id,
			"page": ctx.Query("page"),
		}).WithError(err).Error("Failed to convert page to uint")
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parseing error",
			Description: "We encountered an error while trying to read page from URL",
		})
		return
	}
	logs, err := FetchLogs(uint(id), int(page))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":   id,
			"page": page,
		}).WithError(err).Error("Failed to fetch logs for business")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch the business logs",
		})
		return
	}
	ctx.JSON(200, logs)
}

func (BR *BusinessRouter) FetchEmployees(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithField("id", id).WithError(err).Error("Failed to convert id to uint")
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	list, err := FetchEmployees(uint(id))
	if err != nil {
		logrus.WithField("id", id).WithError(err).Error("Failed to fetch employees for business")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to fetch all employees for the selected business",
		})
		return
	}
	ctx.JSON(200, list)
}

func (BR *BusinessRouter) DeleteBusiness(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithField("id", id).WithError(err).WithError(err).Error("Failed to convert id to uint")
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}

	ctxUserInfo, exists := ctx.Get("userInfo")
	authInfo := ctxUserInfo.(*authinfo.AuthInfo)
	if !exists {
		logrus.Error("Failed to get userinfo in request trying to make an API key")
		ctx.JSON(403, errors.Unauthorized)
		return
	}

	err = DeleteBusiness(authInfo.ID, uint(id))
	if err != nil {
		logrus.WithField("id", id).WithError(err).Error("Failed to delete business")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to delete the given business",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}

func (BR *BusinessRouter) ChangeOwner(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("Failed to convert id to uint")
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}

	reqBody := struct {
		Cid int `json:"cid"`
	}{}
	err = ctx.BindJSON(&reqBody)
	if err != nil {
		logrus.WithField("id", id).WithError(err).Error("Failed to bind body to struct when updating business owner")
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to parse the request data",
		})
		return
	}

	ctxUserInfo, exists := ctx.Get("userInfo")
	authInfo := ctxUserInfo.(*authinfo.AuthInfo)
	if !exists {
		logrus.Error("Failed to get userinfo in request trying to make an API key")
		ctx.JSON(403, errors.Unauthorized)
		return
	}
	err = ChangeOwner(authInfo.ID, uint(id), uint(reqBody.Cid))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":       id,
			"newOwner": reqBody.Cid,
		}).WithError(err).Error("Failed to update the business owner")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an error while trying to update the business owner",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}

func (BR *BusinessRouter) CreateBusiness(ctx *gin.Context) {
	body := struct {
		Label    string `json:"label"`
		Name     string `json:"name"`
		TypeName string `json:"typeName"`
		Owner    uint   `json:"owner"`
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
	ai, err := api.CfxApi.DoRequest("POST", "/business/actions/new", &body, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to create business on cfx")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to create a business",
		})
		return
	}
	if ai.Message != "" {
		logrus.WithField("error", ai.Message).Error("Failed to create business on cfx")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to create a business",
		})
		return
	}
	graylogger.Log("dev:actions:createBusiness", fmt.Sprintf("%d (%s) heeft een bedrijf (%s) aangemaakt", userInfo.ID, userInfo.Username, body.Name), "info", body)
	ctx.JSON(200, gin.H{})
}
