package cfx_models

type CharacterReputation struct {
	Citizenid        uint      `json:"citizenid" gorm:"primaryKey"`
	Crafting         int       `json:"crafting"`
	AmmoCrafting     int       `json:"ammo_crafting" gorm:"column:ammo_crafting"`
	MechanicCrafting int       `json:"mechanic_crafting" gorm:"column:mechanic_crafting"`
	Character        Character `json:"-" gorm:"foreignKey:Citizenid;references:Citizenid"`
}

func (CharacterReputation) TableName() string {
	return "character_reputations"
}
