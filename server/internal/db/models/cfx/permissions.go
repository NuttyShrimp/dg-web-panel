package cfx_models

type Permissions struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	UserSteamId string `json:"steamId" gorm:"index"`
	Users       User   `json:"foreignKey:UserSteamId"`
}
