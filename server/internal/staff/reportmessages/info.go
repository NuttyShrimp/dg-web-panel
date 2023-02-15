package reportmessages

import (
	"degrens/panel/internal/cfx"
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/models"
	"errors"
	"fmt"
)

func SeedReportMessageMember(msg *panel_models.ReportMessage) error {
	if msg == nil {
		return errors.New("Tried to access empty message")
	}
	var messageSender *panel_models.ReportMessageSender
	if msg.MemberID != nil {
		var member panel_models.ReportMember
		db.MariaDB.Client.Where("id = ? AND report_id = ?", msg.MemberID, msg.ReportID).First(&member)
		if member.ID == 0 {
			return fmt.Errorf("could not find a member for %d in report %d", msg.MemberID, msg.ReportID)
		}
		memberInfo, err := cfx.GetCfxUserInfo(member.SteamID)
		if err != nil {
			return err
		}
		messageSender = &panel_models.ReportMessageSender{
			UserInfo: *memberInfo,
			SteamId:  member.SteamID,
		}
	} else if msg.UserID != nil {
		user := db.MariaDB.Repository.GetUserById(*msg.UserID)
		if user.ID == 0 {
			return fmt.Errorf("failed to retrieve user with id %d", msg.UserID)
		}
		steamId := cfx.GetSteamIdFromDiscordId(user.DiscordID)
		messageSender = &panel_models.ReportMessageSender{
			UserInfo: models.UserInfo{
				Username:  user.Username,
				AvatarUrl: user.AvatarUrl,
				Roles:     user.GetRoleNames(),
			},
			SteamId: steamId,
		}
	} else {
		return fmt.Errorf("could not find a valid member for report message with id %d", msg.ID)
	}
	msg.Sender = *messageSender
	return nil
}
