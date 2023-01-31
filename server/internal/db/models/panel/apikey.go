package panel_models

import "time"

type APIKey struct {
	ApiKey    string    `gorm:"primaryKey" json:"key"`
	CreatedAt time.Time `json:"createdAt"`
	Comment   string    `json:"comment"`
	Expiry    time.Time `json:"expiry"`
	UserID    uint      `json:"userId"`
	User      User
}

func (ap *APIKey) Expired() bool {
	return ap.Expiry.Before(time.Now())
}
