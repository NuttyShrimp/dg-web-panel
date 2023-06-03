package bank

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"degrens/panel/lib/graylogger"
	"degrens/panel/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Bank struct {
	*cfx_models.BankAccount
}

func GetBankAccount(accountId string) (*Bank, error) {
	var dbAcc cfx_models.BankAccount
	if err := db.CfxMariaDB.Client.Where(&cfx_models.BankAccount{AccountId: accountId}).First(&dbAcc).Error; err != nil {
		return nil, err
	}
	return &Bank{
		BankAccount: &dbAcc,
	}, nil
}

func (b *Bank) ChangeBalance(userInfo *models.UserInfo, balance float64) error {
	graylogger.Log("admin:bank:updateBalance", fmt.Sprintf("%s(%d) heeft het balans van bankaccount %s (%s) veranded %f -> %f", userInfo.Username, userInfo.ID, b.Name, b.AccountId, b.Balance, balance), "accountId", b.AccountId, "oldBalance", b.Balance, "newBalance", balance)
	if err := db.CfxMariaDB.Client.Model(&cfx_models.BankAccount{}).Where(&cfx_models.BankAccount{AccountId: b.AccountId}).Update("balance", balance).Error; err != nil {
		return err
	}
	ai, err := api.CfxApi.DoRequest("PATCH", "/financials/updateBalance", gin.H{
		"accountId": b.AccountId,
		"balance":   balance,
	}, nil)
	if ai.Message != "" && ai.Response.StatusCode != 401 {
		return errors.New(ai.Message)
	}
	if err != nil {
		return err
	}
	b.Balance = balance
	return nil
}
