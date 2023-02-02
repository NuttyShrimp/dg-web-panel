package cfx_models

import "time"

type PenaltyType string

var (
	BanPenalty  PenaltyType = "ban"
	KickPenalty PenaltyType = "kick"
	WarnPenalty PenaltyType = "warn"
)

type User struct {
	SteamId      string    `json:"steamId,omitempty" gorm:"primaryKey;column:steamid"`
	Name         string    `json:"name,omitempty"`
	License      string    `json:"license,omitempty"`
	Discord      string    `json:"discord,omitempty"`
	Last_Updated time.Time `json:"last_updated,omitempty"`
	Created_At   time.Time `json:"created_at,omitempty"`
}

type Penalties struct {
	Id        uint        `json:"id" gorm:"primaryKey"`
	SteamId   string      `json:"steamId" gorm:"column:steamId"`
	Penalty   PenaltyType `json:"penalty"`
	Reason    string      `json:"reason"`
	Points    uint        `json:"points"`
	Length    int         `json:"length"`
	Date      time.Time   `json:"date"`
	Automated bool        `json:"automated"`
}
