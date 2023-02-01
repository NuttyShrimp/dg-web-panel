package cfx_models

type BusinessType struct {
	Id   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type Business struct {
	Id            uint         `gorm:"primaryKey" json:"id"`
	Label         string       `json:"label"`
	Name          string       `json:"name"`
	BType         uint         `gorm:"column:business_type" json:"-"`
	BankAccountId string       `json:"bankAccountId"`
	BusinessType  BusinessType `gorm:"foreignKey:BType;references:id" json:"type"`
}

type BusinessRole struct {
	Id          uint     `gorm:"primaryKey" json:"id"`
	Name        string   `json:"name"`
	PermMask    uint     `gorm:"column:permissions" json:"-"`
	Permissions []string `json:"permissions" gorm:"-"`
	BusinessId  uint     `json:"-"`
	Business    Business `gorm:"foreignKey:BusinessId;references:id" json:"-"`
}

type BusinessEmployee struct {
	Id         uint         `gorm:"primaryKey" json:"id"`
	IsOwner    bool         `json:"isOwner"`
	CitizenId  uint         `json:"-" gorm:"column:citizenid"`
	RoleId     uint         `json:"-"`
	BusinessId uint         `json:"-"`
	Char       Character    `gorm:"foreignKey:CitizenId;references:Citizenid" json:"character"`
	Role       BusinessRole `gorm:"foreignKey:RoleId;references:id" json:"role"`
	Business   Business     `gorm:"foreignKey:BusinessId;references:id" json:"-"`
}

type BusinessLog struct {
	Id         uint      `json:"id"`
	Type       string    `json:"type"`
	Action     string    `json:"action"`
	BusinessId uint      `json:"businessId"`
	CitizenId  uint      `json:"-" gorm:"column:citizenid"`
	Char       Character `gorm:"foreignKey:CitizenId;references:citizenid" json:"character"`
	Business   Business  `gorm:"foreignKey:BusinessId;references:id" json:"-"`
}

func (BusinessType) TableName() string {
	return "business_type"
}

func (Business) TableName() string {
	return "business"
}

func (BusinessRole) TableName() string {
	return "business_role"
}

func (BusinessEmployee) TableName() string {
	return "business_employee"
}

func (BusinessLog) TableName() string {
	return "business_log"
}
