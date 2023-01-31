package reportmessages

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/cfx"
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"encoding/json"
)

func saveMessage(reportId uint, message interface{}, sender *authinfo.AuthInfo) (*panel_models.ReportMessage, error) {
	msgStr, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	reportMessage := panel_models.ReportMessage{
		Message:  string(msgStr),
		Type:     panel_models.Text,
		ReportID: reportId,
	}
	if sender.AuthMethod == authinfo.CFXToken {
		// Member
		userInfo, err := cfx.GetCfxUserFromId(sender.ID)
		if err != nil {
			return nil, err
		}
		reportMember := db.MariaDB.Repository.GetReportMemberBySteamId(userInfo.SteamId, reportId)
		reportMessage.MemberID = &reportMember.ID
	} else {
		// User
		reportMessage.UserID = &sender.ID
	}
	err = db.MariaDB.Client.Create(&reportMessage).Error
	if err != nil {
		return nil, err
	}
	err = SeedReportMessageMember(&reportMessage)
	if err != nil {
		return nil, err
	}
	return &reportMessage, nil
}
