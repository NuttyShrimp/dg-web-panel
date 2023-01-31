package mariadb

import (
	panel_models "degrens/panel/internal/db/models/panel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	Client *gorm.DB
}

func newRepository(db *gorm.DB) *Repository {
	return &Repository{
		Client: db,
	}
}

func (r *Repository) GetUserById(id uint) panel_models.User {
	var user panel_models.User
	r.Client.Preload(clause.Associations).First(&user, id)
	return user
}

func (r *Repository) GetReportMemberBySteamId(steamId string, reportId uint) panel_models.ReportMember {
	member := panel_models.ReportMember{
		ReportID: reportId,
		SteamID:  steamId,
	}
	r.Client.First(&member)
	return member
}
