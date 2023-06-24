package reports

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/cfxtoken"
	"degrens/panel/internal/cfx"
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/internal/discord"
	"degrens/panel/internal/users"
	"degrens/panel/lib/graylogger"
	"degrens/panel/models"
	"errors"
	"fmt"

	disgo "github.com/disgoorg/disgo/discord"
	"github.com/sirupsen/logrus"

	"github.com/aidenwallis/go-utils/utils"
)

func CreateNewReport(creator, title string, memberIds []string) (uint, error) {
	members := []panel_models.ReportMember{}
	for _, v := range memberIds {
		plyInfo, err := cfx.GetCfxPlayerInfo(v)
		if err != nil {
			return 0, err
		}
		members = append(members, panel_models.ReportMember{
			SteamID: v,
			Name:    plyInfo.Name,
		})
	}
	report := panel_models.Report{
		Title:   title,
		Members: members,
		Creator: creator,
		Open:    true,
	}
	result := db.MariaDB.Client.Create(&report)
	if result.Error != nil {
		return 0, result.Error
	}
	for i := range members {
		members[i].ReportID = report.ID
		db.MariaDB.Client.Save(members[i])
	}
	cfxInput := models.CfxReportAnnouncement{
		ID: report.ID,
		Recvs: utils.SliceMap(members, func(member panel_models.ReportMember) string {
			return member.SteamID
		}),
	}
	ai, err := api.CfxApi.DoRequest("POST", "/admin/report/announce", &cfxInput, nil)
	if err != nil {
		graylogger.Log("reports:announce:error", "Failed to announce a new report to the cfx server", "reportId", report.ID, "title", report.Title)
		logrus.WithField("type", "error").WithError(err).Error("Failed to announce new report")
	}
	if ai.Message != "" {
		graylogger.Log("reports:announce:error", "Failed to announce a new report to the cfx server", "reportId", report.ID, "title", report.Title)
		logrus.WithField("type", "message").WithError(errors.New(ai.Message)).Error("Failed to announce new report")
	}
	graylogger.Log("reports:created", fmt.Sprintf("%s has created a new report with title: %s", creator, title), "members", memberIds)
	err = discord.SendToReportWebHook(&disgo.WebhookMessageCreate{
		Username: "Panel",
		Embeds: []disgo.Embed{
			{
				Title:       "New report",
				Type:        disgo.EmbedTypeRich,
				Description: fmt.Sprintf("%s has created a new report with title: %s", creator, title),
				Fields: []disgo.EmbedField{
					{
						Name:  "URL",
						Value: fmt.Sprintf("https://panel.degrensrp.be/staff/reports/%d", report.ID),
					},
				},
			},
		},
	})
	return report.ID, err
}

func AddMemberToReport(userId string, reportId uint, steamId string) error {
	report := panel_models.Report{}
	err := db.MariaDB.Client.First(&report, reportId).Error
	if err != nil {
		return err
	}
	if report.ID == 0 {
		return fmt.Errorf("failed to find report with id %d while adding new member", reportId)
	}
	member := panel_models.ReportMember{
		SteamID:  steamId,
		ReportID: report.ID,
	}
	plyInfo, err := cfx.GetCfxPlayerInfo(steamId)
	if err != nil {
		member.Name = plyInfo.Name
	}
	graylogger.Log("reports:member_add", fmt.Sprintf("%s has added a new member(%s) to report %d", userId, steamId, reportId), "steamId", steamId, "reportId", reportId)
	return db.MariaDB.Client.Create(&member).Error
}

func FetchReports(titleFilter string, offset int, includeOpen, includeClosed bool, authInfo *authinfo.AuthInfo) (*[]panel_models.Report, error) {
	dbQuery := db.MariaDB.Client.Model(&panel_models.Report{})
	if authInfo == nil {
		return nil, errors.New("Failed to get authentication info")
	}
	if authInfo.AuthMethod == authinfo.CFXToken && !users.DoesUserHaveRole(authInfo.Roles, "staff") {
		tokenInfo := cfxtoken.GetInfoForToken(authInfo.ID)
		if tokenInfo == nil {
			return nil, errors.New("Failed to get info bound to cfx token")
		}
		dbQuery = dbQuery.Joins("JOIN report_members ON reports.id = report_members.report_id AND report_members.steam_id = ?", tokenInfo.SteamId)
	} else {
		dbQuery = dbQuery.Preload("Members")
	}
	dbQuery = dbQuery.Offset(offset*25).Limit(25).Order("id desc").Where("title LIKE ?", "%"+titleFilter+"%")
	if includeOpen && !includeClosed {
		dbQuery = dbQuery.Where("open = ?", true)
	} else if includeClosed && !includeOpen {
		dbQuery = dbQuery.Where("open = ?", false)
	}

	dbReports := []panel_models.Report{}

	err := dbQuery.Find(&dbReports).Error
	return &dbReports, err
}

func FetchReport(reportId uint) (*panel_models.Report, error) {
	report := panel_models.Report{}
	err := db.MariaDB.Client.Preload("Members").First(&report, reportId).Error
	return &report, err
}

func FetchReportCount(titleFilter string, offset int, includeOpen, includeClosed bool) (int64, error) {
	var reportCount int64
	dbQuery := db.MariaDB.Client.Model(&panel_models.Report{}).Where("title LIKE ?", "%"+titleFilter+"%").Select("id")
	if includeOpen && !includeClosed {
		dbQuery = dbQuery.Where("open = ?", true)
	} else if includeClosed && !includeOpen {
		dbQuery = dbQuery.Where("open = ?", false)
	}
	err := dbQuery.Count(&reportCount).Error
	return reportCount, err
}
