package panel_models

type Role struct {
	Name   string `gorm:"primaryKey"`
	UserId uint   `gorm:"primaryKey"`
}
