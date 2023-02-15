package reportmessages

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/cfxtoken"
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"encoding/json"
	"errors"
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
		tokenInfo := cfxtoken.GetInfoForToken(sender.ID)
		if tokenInfo == nil {
			return nil, errors.New("Failed to get info bound to cfx token")
		}
		reportMember := db.MariaDB.Repository.GetReportMemberBySteamId(tokenInfo.SteamId, reportId)
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
