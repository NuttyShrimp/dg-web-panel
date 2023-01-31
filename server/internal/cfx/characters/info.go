package characters

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
)

func GetCharacterInfo(cid uint) (cfx_models.Character, error) {
	char := cfx_models.Character{}
	err := db.CfxMariaDB.Client.Model(&char).Preload("Data").Preload("Info").Where("citizenid = ?", cid).First(&char).Error
	return char, err
}
