package cfx_models

type RealEstateLocation struct {
	ID       uint                       `gorm:"primaryKey" json:"id"`
	Name     string                     `gorm:"column:name" json:"name"`
	Garage   string                     `gorm:"column:garage" json:"garage"`
	Clothing string                     `gorm:"column:clothing" json:"clothing"`
	Logout   string                     `gorm:"column:logout" json:"logout"`
	Stash    string                     `gorm:"column:stash" json:"stash"`
	Access   []RealEstateLocationAccess `gorm:"foreignKey:LocationId;references:id" json:"access"`
}

type RealEstateLocationAccess struct {
	LocationId uint      `gorm:"primaryKey" json:"locationId"`
	Owner      bool      `gorm:"column:owner" json:"owner"`
	CitizenId  uint      `gorm:"column:cid" json:"citizenId"`
	Character  Character `gorm:"foreignKey:CitizenId;references:citizenid" json:"character"`
}

func (*RealEstateLocation) TableName() string {
	return "realestate_locations"
}

func (*RealEstateLocationAccess) TableName() string {
	return "realestate_location_access"
}
