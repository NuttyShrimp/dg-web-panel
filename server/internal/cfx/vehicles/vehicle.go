package vehicles

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
)

func FetchForCid(cid uint) (*[]cfx_models.PlayerVehicles, error) {
	vehs := []cfx_models.PlayerVehicles{}
	err := db.CfxMariaDB.Client.Where("cid = ?", cid).Find(&vehs).Error
	return &vehs, err
}
