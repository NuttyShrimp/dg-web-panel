package characters

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"degrens/panel/lib/log"
	"degrens/panel/models"
	"time"

	"github.com/aidenwallis/go-utils/utils"
)

var cidCache = []models.CfxCharacter{}
var cacheTimer *time.Timer
var logger *log.Logger

func DoesCIDExist(cid uint) bool {
	chars := int64(0)
	db.CfxMariaDB.Client.Where("citizenid = ?", cid).Model(&cfx_models.Character{}).Count(&chars)
	return chars > 0
}

func FetchAllCids() ([]uint, error) {
	cids := []cfx_models.Character{}
	err := db.CfxMariaDB.Client.Find(&cids).Select("citizenid").Error
	return utils.SliceMap(cids, (func(c cfx_models.Character) uint { return c.Citizenid })), err
}
