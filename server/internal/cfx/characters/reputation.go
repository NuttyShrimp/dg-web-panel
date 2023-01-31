package characters

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"

	"gorm.io/gorm"
)

func GetCharacterReputation(cid uint) (*cfx_models.CharacterReputation, error) {
	rep := cfx_models.CharacterReputation{}
	err := db.CfxMariaDB.Client.Model(&rep).Select("crafting, ammo_crafting, mechanic_crafting").Where("citizenid = ?", cid).First(&rep).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return &rep, err
}
