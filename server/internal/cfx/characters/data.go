package characters

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
)

func GetCharacterDataForCid(cid uint) (*cfx_models.Character, error) {
	data := cfx_models.Character{}
	err := db.CfxMariaDB.Client.Where("citizenid = ?", cid).First(&data).Error
	return &data, err
}

func GetCharactersForSteamId(steamId string) ([]*cfx_models.Character, error) {
	chars := []*cfx_models.Character{}
	err := db.CfxMariaDB.Client.Preload("Info").Where(&cfx_models.Character{UserSteamId: steamId}).Find(&chars).Error
	return chars, err
}

func GetAllCharacters() (*[]cfx_models.Character, error) {
	data := []cfx_models.Character{}
	err := db.CfxMariaDB.Client.Preload("User").Preload("Data").Preload("Info").Find(&data).Error
	return &data, err
}
