package cfx_models

type FlyerRequest struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	CitizenId uint      `json:"citizenid" gorm:"column:cid"`
	Link      string    `json:"link"`
	Approved  bool      `json:"approved"`
	Character Character `gorm:"foreignKey:CitizenId" json:"character"`
}

func (FlyerRequest) TableName() string {
	return "flyer_request"
}
