package reports

import (
	"degrens/panel/internal/auth/authinfo"
	"degrens/panel/internal/auth/cfxtoken"
	"degrens/panel/internal/cfx"
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/lib/graylogger"
	"errors"
	"fmt"
)

func CreateNewReport(creator, title string, memberIds, tagNames []string) (uint, error) {
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
	tags := []panel_models.ReportTag{}
	for _, tagName := range tagNames {
		tag := panel_models.ReportTag{
			Name: tagName,
		}
		db.MariaDB.Client.First(&tag)
		tags = append(tags, tag)
	}
	report := panel_models.Report{
		Title:   title,
		Tags:    tags,
		Members: members,
		Creator: creator,
		Open:    true,
	}
	result := db.MariaDB.Client.Create(&report)
	if result.Error != nil {
		return 0, result.Error
	}
	for _, member := range members {
		member.ReportID = report.ID
		db.MariaDB.Client.Save(&member)
	}
	graylogger.Log("reports:created", fmt.Sprintf("%s has created a new report with title: %s", title), "members", memberIds, "tags", tagNames)
	return report.ID, nil
}

func AddMemberToReport(userId string, reportId uint, steamId string) error {
	report := panel_models.Report{}
	err := db.MariaDB.Client.First(&report, reportId).Error
	if err != nil {
		return err
	}
	if report.ID == 0 {
		return errors.New(fmt.Sprintf("failed to find report with id %d while adding new member", reportId))
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

func FetchReports(titleFilter string, offset int, tags []string, includeOpen, includeClosed bool, authInfo *authinfo.AuthInfo) (*[]panel_models.Report, error) {
	dbQuery := db.MariaDB.Client.Preload("Tags")
	if len(tags) > 0 {
		dbQuery.Joins("inner join report_tags_link rtl on rtl.report_id = reports.id ").
			Joins("inner join report_tags t on t.name = rtl.report_tag_name").
			Where("t.name IN ?", tags)
	}
	if authInfo != nil && authInfo.AuthMethod == authinfo.CFXToken {
		tokenInfo := cfxtoken.GetInfoForToken(authInfo.ID)
		if tokenInfo == nil {
			return nil, errors.New("Failed to get info bound to cfx token")
		}
		dbQuery.Joins("Members", db.MariaDB.Client.Where(&panel_models.ReportMember{SteamID: tokenInfo.SteamId}))
	} else {
		dbQuery.Preload("Members")
	}
	dbQuery.Offset(offset*25).Limit(25).Order("id desc").Where("title LIKE ?", "%"+titleFilter+"%")
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
	err := db.MariaDB.Client.Preload("Tags").Preload("Members").First(&report, reportId).Error
	return &report, err
}

func FetchReportCount(titleFilter string, offset int, tags []string, includeOpen, includeClosed bool) (int64, error) {
	var reportCount int64
	dbQuery := db.MariaDB.Client.Raw("SELECT COUNT(id) FROM reports WHERE (SELECT COUNT(DISTINCT report_id) FROM report_tags_link WHERE report_tag_name IN ('cry-baby', 'bug-report') AND report_id = reports.id) > 0 AND reports.title LIKE ?", "%"+titleFilter+"%")
	if includeOpen && !includeClosed {
		dbQuery = dbQuery.Where("open = ?", true)
	} else if includeClosed && !includeOpen {
		dbQuery = dbQuery.Where("open = ?", false)
	}
	err := dbQuery.Count(&reportCount).Error
	return reportCount, err
}
