package cfx_models

import "time"

type PenaltyType string

var (
	BanPenalty  PenaltyType = "ban"
	KickPenalty PenaltyType = "kick"
	WarnPenalty PenaltyType = "warn"
)

type User struct {
	SteamId      string      `json:"steamId,omitempty" gorm:"primaryKey;column:steamid"`
	Name         string      `json:"name,omitempty"`
	License      string      `json:"license,omitempty"`
	DiscordId    string      `json:"discordId,omitempty" gorm:"column:discord"`
	Last_Updated time.Time   `json:"last_updated,omitempty"`
	Created_At   time.Time   `json:"created_at,omitempty"`
	Points       AdminPoints `gorm:"foreignKey:SteamId" json:"points"`
}

// nolint:gocritic,unused
func (item *User) compare(other User) int {
	if item.SteamId == other.SteamId {
		return 0
	}
	if item.SteamId < other.SteamId {
		return -1
	}
	return 1
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
