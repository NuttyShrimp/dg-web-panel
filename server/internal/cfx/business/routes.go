package business

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/routes"
	"degrens/panel/lib/errors"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var br BusinessRouter

type BusinessRouter struct {
	routes.Router
}

func NewBusinessRouter(rg *gin.RouterGroup, logger log.Logger) {
	br = BusinessRouter{
		Router: routes.Router{
			RouterGroup: rg.Group("/business"),
			Logger:      logger,
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
			BR.Logger.Error("Failed to convert cid to uint", "error", err, "cid", cid64)
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
		BR.Logger.Error("Failed to fetch all businesses", "error", err)
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
		BR.Logger.Error("Failed to convert id to uint", "error", err, "id", ctx.Param("id"))
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	total, err := FetchLogCount(uint(id))
	if err != nil {
		BR.Logger.Error("Failed to fetch total amount of logs for a business", "error", err, "business_id", id)
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
		BR.Logger.Error("Failed to convert id to uint", "error", err, "id", id)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	page, err := strconv.ParseUint(ctx.DefaultQuery("page", "0"), 10, 32)
	if err != nil {
		BR.Logger.Error("Failed to convert page to uint", "error", err, "id", id, "page", ctx.Query("page"))
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parseing error",
			Description: "We encountered an error while trying to read page from URL",
		})
		return
	}
	logs, err := FetchLogs(uint(id), int(page))
	if err != nil {
		BR.Logger.Error("Failed to fetch logs for business", "error", err, "id", id, "page", page)
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
		BR.Logger.Error("Failed to convert id to uint", "error", err, "id", id)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}
	list, err := FetchEmployees(uint(id))
	if err != nil {
		BR.Logger.Error("Failed to fetch employees for business", "error", err, "id", id)
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
		BR.Logger.Error("Failed to convert id to uint", "error", err, "id", id)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to convert the given Business id to a valid number",
		})
		return
	}

	ctxUserInfo, exists := ctx.Get("userInfo")
	authInfo := ctxUserInfo.(*authinfo.AuthInfo)
	if !exists {
		BR.Logger.Error("Failed to get userinfo in request trying to make an API key")
		ctx.JSON(403, errors.Unauthorized)
		return
	}

	err = DeleteBusiness(authInfo.ID, uint(id))
	if err != nil {
		BR.Logger.Error("Failed to delete business", "error", err, "id", id)
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
		BR.Logger.Error("Failed to convert id to uint", "error", err, "id", id)
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
		BR.Logger.Error("Failed to bind body to struct when updating business owner", "error", err, "id", id)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to parse the request data",
		})
		return
	}

	ctxUserInfo, exists := ctx.Get("userInfo")
	authInfo := ctxUserInfo.(*authinfo.AuthInfo)
	if !exists {
		BR.Logger.Error("Failed to get userinfo in request trying to make an API key")
		ctx.JSON(403, errors.Unauthorized)
		return
	}
	err = ChangeOwner(authInfo.ID, uint(id), uint(reqBody.Cid))
	if err != nil {
		BR.Logger.Error("Failed to update the business owner", "error", err, "id", id, "newOwner", reqBody.Cid)
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
		BR.Logger.Error("Failed to bind body to struct when createing business", "error", err)
		ctx.JSON(400, models.RouteErrorMessage{
			Title:       "Request error",
			Description: "We encountered an issue while reading the info from your request",
		})
		return
	}
	ai, err := api.CfxApi.DoRequest("POST", "/business/actions/new", &body, nil)
	if err != nil {
		BR.Logger.Error("Failed to create business on cfx", "error", err)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to create a business",
		})
		return
	}
	if ai.Message != "" {
		BR.Logger.Error("Failed to create business on cfx", "error", ai.Message)
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: "We encountered an issue while trying to create a business",
		})
		return
	}
	ctx.JSON(200, gin.H{})
}
