package bank

import (
	cfx_models "degrens/panel/internal/db/models/cfx"
)

func getPermissions(permNum uint) *cfx_models.BankPermissions {
	perms := cfx_models.BankPermissions{
		Deposit:      permNum&1 == 1,
		Withdraw:     permNum&2 == 2,
		Transfer:     permNum&4 == 4,
		Transactions: permNum&8 == 8,
		Owner:        permNum&16 == 16,
	}
	return &perms
}
