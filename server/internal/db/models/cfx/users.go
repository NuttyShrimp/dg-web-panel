package cfx_models

import "time"

type User struct {
	SteamId      string    `json:"steamId,omitempty" gorm:"primaryKey;column:steamid"`
	Name         string    `json:"name,omitempty"`
	License      string    `json:"license,omitempty"`
	Discord      string    `json:"discord,omitempty"`
	Last_Updated time.Time `json:"last_updated,omitempty"`
	Created_At   time.Time `json:"created_at,omitempty"`
}
