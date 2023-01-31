package panel_models

type Notes struct {
	BaseModel
	CreatorID uint   `json:"-"`
	Note      string `json:"note"`
	User      User   `gorm:"foreignKey:CreatorID;references:id" json:"user"`
}
