package bank

import (
	"degrens/panel/internal/routes"
	"degrens/panel/internal/users"
	"degrens/panel/lib/errors"
	"degrens/panel/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BankRouter struct {
	routes.Router
}

func NewBankRouter(rg *gin.RouterGroup) {
	br := BankRouter{
		Router: routes.Router{
			RouterGroup: rg.Group("/bank"),
		},
	}
	br.RegisterRoutes()
}

func (BR *BankRouter) RegisterRoutes() {
	BR.RouterGroup.GET("/:cid", BR.FetchAccounts)
	BR.RouterGroup.PATCH("/:id/balance", BR.UpdateBalance)
}

func (BR *BankRouter) FetchAccounts(ctx *gin.Context) {
	cid, err := strconv.ParseInt(ctx.Param("cid"), 10, 32)
	if err != nil {
		logrus.WithField("cid", ctx.Param("cid")).WithError(err).Error("Failed to convert citizenid to uint")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the character you are trying to fetch the bank accounts for",
		})
		return
	}

	accs, err := getAccounts(uint(cid))
	if err != nil {
		logrus.WithField("cid", ctx.Param("cid")).WithError(err).Error("Failed to fetch bank accounts")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: fmt.Sprintf("We encountered an error while trying to fetch the bank accounts for cid: %d", cid),
		})
		return
	}

	ctx.JSON(200, accs)
}

func (BR *BankRouter) UpdateBalance(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while trying to identify the bank account you are trying to fetch the bank accounts for",
		})
		return
	}

	var body struct {
		Balance float64 `json:"balance"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		logrus.WithField("id", id).WithError(err).Error("Failed to bind request body")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Parsing error",
			Description: "We encountered an error while getting the request information",
		})
		return
	}

	bank, err := GetBankAccount(id)
	if err != nil {
		logrus.WithField("accountId", id).WithError(err).Error("Failed to fetch bank account")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: fmt.Sprintf("We encountered an error while trying to fetch the bank account for id: %s", id),
		})
		return
	}

	userInfo, err := users.GetUserInfo(ctx)
	if err != nil {
		logrus.Error("Failed to get userinfo while giving vehicle to player")
		ctx.JSON(403, errors.Unauthorized)
		return
	}

	if err = bank.ChangeBalance(userInfo, body.Balance); err != nil {
		logrus.WithField("id", id).WithError(err).Error("Failed to update bank account balance")
		ctx.JSON(500, models.RouteErrorMessage{
			Title:       "Server error",
			Description: fmt.Sprintf("We encountered an error while updating bank account (%s) to balance %f", id, body.Balance),
		})
		return
	}

	ctx.JSON(200, gin.H{})
}
