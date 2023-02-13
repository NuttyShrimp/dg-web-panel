package discord

import (
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

// Refreshes the data of an existing user
func UpdateUserInfo(identity DiscordIdentity) *panel_models.User {
	var user panel_models.User
	db.MariaDB.Client.FirstOrCreate(&user, panel_models.User{
		DiscordID: identity.Id,
	})
	user.Username = identity.Username
	user.DiscordID = identity.Id
	user.AvatarUrl = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s", identity.Id, identity.AvatarHash)
	user.Roles = []panel_models.Role{}
	for _, roleName := range GetRegisterdRolesForIdentity(identity) {
		var role panel_models.Role
		db.MariaDB.Client.FirstOrCreate(&role, panel_models.Role{
			Name:   roleName,
			UserId: user.ID,
		})
		user.Roles = append(user.Roles, role)
	}
	db.MariaDB.Client.Save(&user)
	return &user
}

// Checks if an identity has atleast 1 valid role
func GetRegisterdRolesForIdentity(identity DiscordIdentity) []string {
	registerdRoles := []string{}
	for _, userRole := range identity.Roles {
		for _, srvRole := range info.Roles {
			if userRole == srvRole.Id {
				registerdRoles = append(registerdRoles, srvRole.Name)
			}
		}
	}
	return registerdRoles
}

func AssignTokenToUser(userId uint, token *oauth2.Token) {
	RemoveUserTokens(userId)

	DBToken := panel_models.DiscordTokens{
		UserID:       userId,
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	db.MariaDB.Client.Create(&DBToken)
}

func RemoveUserTokens(userId uint) {
	Token := panel_models.DiscordTokens{
		UserID: userId,
	}
	oldTokens := []panel_models.DiscordTokens{}
	db.MariaDB.Client.Where(&Token).Find(&oldTokens)
	for i := range oldTokens {
		err := RevokeAuthToken(oldTokens[i].Token)
		if err != nil {
			logger.Error("Failed to remove user token", "error", err)
		}
	}
	db.MariaDB.Client.Delete(&panel_models.DiscordTokens{}, "user_id = ?", userId)
}

// Returns if a token is expired
func IsTokenExpired(userId uint) bool {
	var DBToken panel_models.DiscordTokens
	db.MariaDB.Client.First(&DBToken, "user_id = ?", userId)
	return DBToken.Expiry.Before(time.Now())
}
