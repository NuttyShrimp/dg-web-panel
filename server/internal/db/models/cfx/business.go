package cfx_models

type BusinessType struct {
	Id   uint `gorm:"primaryKey"`
	Name string
}

type Business struct {
	Id            uint `gorm:"primaryKey"`
	Label         string
	Name          string
	BType         uint `gorm:"column:business_type"`
	BankAccountId string
	BusinessType  BusinessType `gorm:"foreignKey:BType;references:id"`
}

type BusinessRole struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Permissions uint
	BusinessId  uint
	Business    Business `gorm:"foreignKey:BusinessId;references:id"`
}

type BusinessEmployee struct {
	Id         uint `gorm:"primaryKey"`
	IsOwner    bool
	CitizenId  uint
	RoleId     uint
	BusinessId uint
	Char       Character    `gorm:"foreignKey:CitizenId;references:citizenid"`
	Role       BusinessRole `gorm:"foreignKey:RoleId;references:id"`
	Business   Business     `gorm:"foreignKey:BusinessId;references:id"`
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
