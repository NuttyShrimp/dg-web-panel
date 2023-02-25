package cfx_models

import "time"

type AdminPoints struct {
	SteamId   string    `json:"steamId" gorm:"index;column:steamid"`
	Points    int       `json:"points"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (AdminPoints) TableName() string {
	return "admin_points"
}
