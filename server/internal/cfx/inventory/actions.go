package inventory

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"

	"gorm.io/datatypes"
)

func FetchItemsByMetadata(metadata map[string]string) ([]cfx_models.InventoryItem, error) {
	items := []cfx_models.InventoryItem{}
	type_query := datatypes.JSONQuery("metadata")
	for k, v := range metadata {
		if v == "" {
			type_query.HasKey(k)
		} else {
			type_query.Equals(v, k)
		}
	}
	err := db.CfxMariaDB.Client.Find(&items, type_query).Error
	return items, err
}
