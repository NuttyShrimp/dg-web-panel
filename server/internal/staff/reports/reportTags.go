package reports

import (
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"fmt"
)

func GetTags() (*[]panel_models.ReportTag, error) {
	tags := []panel_models.ReportTag{}
	result := db.MariaDB.Client.Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tags, nil
}

func NewReportTag(name, color string) error {
	tag := panel_models.ReportTag{
		Name:  name,
		Color: color,
	}
	result := db.MariaDB.Client.Create(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateReportTags(reportId uint, tagNames []string) error {
	tags := []panel_models.ReportTag{}
	for _, tagName := range tagNames {
		tag := panel_models.ReportTag{
			Name: tagName,
		}
		db.MariaDB.Client.First(&tag)
		tags = append(tags, tag)
	}
	var report panel_models.Report
	db.MariaDB.Client.First(&report, reportId)
	if report.ID == 0 {
		return fmt.Errorf("Failed to find report with id %d", reportId)
	}
	report.Tags = tags
	db.MariaDB.Client.Save(&report)
	return nil
}
