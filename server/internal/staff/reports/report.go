package reports

import (
	"degrens/panel/internal/cfx"
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"

	"github.com/aidenwallis/go-utils/utils"
)

type Report struct {
	Data *panel_models.Report
}

func CreateReport(report *panel_models.Report) Report {
	return Report{
		Data: report,
	}
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
