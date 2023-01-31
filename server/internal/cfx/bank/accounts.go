package bank

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
)

func getAccounts(cid uint) (*[]cfx_models.BankAccount, error) {
	acc_perms := []cfx_models.BankAccountAccess{}
	perms_err := db.CfxMariaDB.Client.Where("cid = ?", cid).Find(&acc_perms).Error

	accs := []cfx_models.BankAccount{}
	for _, perm := range acc_perms {
		acc := cfx_models.BankAccount{}
		acc_err := db.CfxMariaDB.Client.Where("account_id = ?", perm.AccountId).Find(&acc).Error
		if acc_err != nil {
			return nil, acc_err
		}
		acc.Permissions = *getPermissions(uint(perm.Access_Level))
		accs = append(accs, acc)
	}
	return &accs, perms_err
}
