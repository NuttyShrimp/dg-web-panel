package cfx_models

import "time"

type Character struct {
	Citizenid    uint          `json:"citizenid" gorm:"primaryKey"`
	Last_Updated time.Time     `json:"last_updated"`
	Created_At   time.Time     `json:"created_at"`
	UserSteamId  string        `json:"steamId" gorm:"column:steamid"`
	User         User          `json:"user" gorm:"foreignKey:UserSteamId;references:steamid"`
	Data         CharacterData `json:"data" gorm:"foreignKey:Citizenid;references:citizenid"`
	Info         CharacterInfo `json:"info" gorm:"foreignKey:Citizenid;references:citizenid"`
}

type CharacterInfo struct {
	Citizenid    uint      `json:"citizenid" gorm:"primaryKey"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Birthdate    string    `json:"birthdate"`
	Gender       int       `json:"gender"`
	Nationality  string    `json:"nationality"`
	Phone        string    `json:"phone"`
	Last_Updated time.Time `json:"last_updated"`
	Created_At   time.Time `json:"created_at"`
}

type CharacterData struct {
	Citizenid    uint      `json:"citizenid" gorm:"primaryKey"`
	Position     string    `json:"position"`
	Metadata     string    `json:"metadata"`
	Last_Updated time.Time `json:"last_updated"`
	Created_At   time.Time `json:"created_at"`
}

func (CharacterInfo) TableName() string {
	return "character_info"
}

func (CharacterData) TableName() string {
	return "character_data"
}
