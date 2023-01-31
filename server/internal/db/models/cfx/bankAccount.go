package cfx_models

import "time"

type AccountType string

var (
	StandardAccount AccountType = "standard"
	SavingsAccount              = "savings"
	BusinessAccount             = "business"
)

type BankAccount struct {
	AccountId   string          `gorm:"primaryKey" json:"accountId"`
	Name        string          `json:"name"`
	Type        AccountType     `json:"accountType"`
	Balance     float64         `json:"balance"`
	Updated_At  time.Time       `json:"updated_at"`
	Permissions BankPermissions `json:"permissions" gorm:"-"`
}

type BankAccountAccess struct {
	AccountId    string `gorm:"primaryKey" json:"accountId"`
	CitizenId    int    `gorm:"primaryKey,column:cid" json:"citizenId"`
	Access_Level int    `json:"access_level"`
}

type BankPermissions struct {
	Deposit      bool `json:"deposit"`
	Withdraw     bool `json:"withdraw"`
	Transfer     bool `json:"transfer"`
	Transactions bool `json:"transactions"`
	Owner        bool `json:"owner"`
}

func (BankAccount) TableName() string {
	return "bank_accounts"
}

func (BankAccountAccess) TableName() string {
	return "bank_accounts_access"
}
