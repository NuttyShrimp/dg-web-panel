package reports

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/cfxtoken"
	"degrens/panel/internal/cfx"
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/internal/users"
	"degrens/panel/lib/graylogger"
	"degrens/panel/models"
	"encoding/json"
	"errors"

	"github.com/aidenwallis/go-utils/utils"
)

type Report struct {
	Data *panel_models.Report
}

var reportCache = make(map[uint]*Report)

func CreateReport(report *panel_models.Report) *Report {
	if r, ok := reportCache[report.ID]; ok {
		return r
	}
	r := Report{
		Data: report,
	}
	reportCache[report.ID] = &r
	return &r
}

func GetReport(id uint) (*Report, error) {
	if r, ok := reportCache[id]; ok {
		return r, nil
	}
	reportData, err := FetchReport(id)
	if err != nil {
		return nil, err
	}
	return CreateReport(reportData), nil
}

func (r *Report) saveToDB() error {
	return db.MariaDB.Client.Save(&r.Data).Error
}

func (r *Report) AddMember(memberSteamId string) error {
	_, ok := utils.SliceFind(r.Data.Members, func(mem panel_models.ReportMember) bool {
		return mem.SteamID == memberSteamId
	})
	if ok {
		return nil
	}
	plyInfo, err := cfx.GetCfxPlayerInfo(memberSteamId)
	if err != nil {
		return err
	}
	member := panel_models.ReportMember{
		ReportID: r.Data.ID,
		SteamID:  memberSteamId,
		Name:     plyInfo.Name,
	}
	err = db.MariaDB.Client.Create(&member).Error
	if err != nil {
		return err
	}
	r.Data.Members = append(r.Data.Members, member)
	return nil
}

func (r *Report) RemoveMember(memberSteamId string) error {
	member, ok := utils.SliceFind(r.Data.Members, func(mem panel_models.ReportMember) bool {
		return mem.SteamID == memberSteamId
	})
	if !ok {
		return nil
	}
	err := db.MariaDB.Client.Delete(&member).Error
	if err != nil {
		return err
	}
	r.Data.Members = utils.SliceFilter(r.Data.Members, func(mem panel_models.ReportMember) bool {
		return mem.SteamID != memberSteamId
	})
	return nil
}

func (r *Report) ToggleState(open bool) error {
	r.Data.Open = open
	return r.saveToDB()
}

func (r *Report) AddMessage(reportId uint, message interface{}, sender *authinfo.AuthInfo) (*panel_models.ReportMessage, error) {
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

	cfxInput := models.CfxReportAnnouncement{
		ID: r.Data.ID,
		Recvs: utils.SliceMap(r.Data.Members, func(member panel_models.ReportMember) string {
			return member.SteamID
		}),
	}
	err = api.CfxApi.Post("/admin/report/announce", &cfxInput, nil)
	if err != nil {
		userId, err := users.GetUserIdentifier(sender)
		if err != nil {
			return nil, err
		}
		graylogger.Log("reports:announce:error", "Failed to announce a new report message to the cfx server", "reportId", r.Data.ID, "message", reportMessage.Message, "sender", userId)
	}

	return &reportMessage, nil
}
