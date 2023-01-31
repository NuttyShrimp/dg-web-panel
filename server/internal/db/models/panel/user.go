package panel_models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DiscordID string `gorm:"unique_index"`
	Username  string
	AvatarUrl string
	Roles     []Role
}

func (u User) GetRoleNames() []string {
	roleNames := []string{}
	for _, role := range u.Roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames
}
