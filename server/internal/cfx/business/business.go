package business

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"

	"gorm.io/gorm"
)

type Filters struct {
	cid *uint
}

func FetchBusinesses(filter *Filters) ([]cfx_models.Business, error) {
	businesses := []cfx_models.Business{}
	dbQuery := db.CfxMariaDB.Client.Preload("BusinessType")
	if filter == nil {
		err := dbQuery.Find(&businesses).Error
		return businesses, err
	}
	if filter.cid != nil && *filter.cid > 0 {
		dbQuery.Joins("BusinessEmployee", db.CfxMariaDB.Client.Where(cfx_models.BusinessEmployee{CitizenId: *filter.cid}))
	}
	err := dbQuery.Find(&businesses).Error
	return businesses, err
}

func FetchLogs(businessId uint, page int) ([]cfx_models.BusinessLog, error) {
	logs := []cfx_models.BusinessLog{}
	err := db.CfxMariaDB.Client.Preload("Char").Preload("Char.User").Preload("Char.Info").Where("business_id = ?", businessId).Offset(page * 50).Limit(50).Find(&logs).Error
	if err == gorm.ErrRecordNotFound {
		return logs, nil
	}
	return logs, err
}

func FetchLogCount(businessId uint) (int64, error) {
	var count int64
	err := db.CfxMariaDB.Client.Model(&cfx_models.BusinessLog{}).Where("business_id = ?", businessId).Count(&count).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return count, err
}

func FetchEmployees(businessId uint) ([]cfx_models.BusinessEmployee, error) {
	employees := []cfx_models.BusinessEmployee{}
	err := db.CfxMariaDB.Client.Preload("Role").Preload("Char").Preload("Char.User").Preload("Char.Info").Where("business_id = ?", businessId).Find(&employees).Error
	for i, employee := range employees {
		employee.Role.Permissions = bitToPermList(employee.Role.PermMask)
		employees[i] = employee
	}
	return employees, err
}
