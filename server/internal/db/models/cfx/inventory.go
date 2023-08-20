package cfx_models

import "gorm.io/datatypes"

type InventoryItem struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name"`
	Inventory   string         `json:"inventory"`
	Position    string         `json:"position"`
	HotKey      int            `json:"hotkey"`
	MetaData    datatypes.JSON `json:"metadata" gorm:"column:metadata"`
	DestroyDate int            `json:"destroyDate"`
	Rotate      int            `json:"-"`
}

func (InventoryItem) TableName() string {
	return "inventory_items"
}
