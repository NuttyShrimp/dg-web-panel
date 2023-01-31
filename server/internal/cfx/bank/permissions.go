package bank

import (
	cfx_models "degrens/panel/internal/db/models/cfx"
)

func getPermissions(permNum uint) *cfx_models.BankPermissions {
	perms := cfx_models.BankPermissions{
		Deposit:      permNum&1 == 0,
		Withdraw:     permNum&2 == 0,
		Transfer:     permNum&4 == 0,
		Transactions: permNum&8 == 0,
		Owner:        permNum&16 == 0,
	}
	return &perms
}
