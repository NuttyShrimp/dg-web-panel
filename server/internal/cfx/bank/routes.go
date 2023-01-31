package bank

import (
	"degrens/panel/internal/routes"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BankRouter struct {
	routes.Router
}

func NewBankRouter(rg *gin.RouterGroup, logger *log.Logger) {
	br := BankRouter{
		Router: routes.Router{
			RouterGroup: rg.Group("/bank"),
			Logger:      *logger,
		},
	}
	br.RegisterRoutes()
}

func (BR *BankRouter) RegisterRoutes() {
	BR.RouterGroup.GET("/:cid", BR.FetchAccounts)
}

func (BR *BankRouter) FetchAccounts(ctx *gin.Context) {
	cid, err := strconv.ParseInt(ctx.Param("cid"), 10, 32)
	if err != nil {
		BR.Logger.Error("Failed to convert citizenid to uint", "error", err, "cid", ctx.Param("cid"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the character you are trying to fetch the bank accounts for",
		})
		return
	}

	accs, err := getAccounts(uint(cid))
	if err != nil {
		BR.Logger.Error("Failed to fetch bank accounts", "error", err, "cid", ctx.Param("cid"))
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: fmt.Sprintf("We encountered an error while trying to fetch the bank accounts for cid: %d", cid),
		})
		return
	}

	ctx.JSON(200, accs)
}
