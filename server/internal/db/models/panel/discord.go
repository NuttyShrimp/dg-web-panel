package panel_models

import "time"

type DiscordTokens struct {
	Token        string `gorm:"primaryKey"`
	RefreshToken string
	Expiry       time.Time
	UserID       uint
	User         User
}
