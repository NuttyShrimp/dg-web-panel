package business

import (
	"degrens/panel/internal/api"
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"
	"degrens/panel/lib/cache"
	"degrens/panel/lib/graylogger"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Filters struct {
	cid *uint
}

// TODO: add ability to empty cache
var businessCache = cache.InitRefreshCache[cfx_models.Business, uint](5*time.Minute, getBusinessById)

// TODO: create overhead Business struct with functions tied to

func getBusinessById(id uint) *cfx_models.Business {
	business := cfx_models.Business{}
	db.CfxMariaDB.Client.Preload("BusinessType").Find(&business, id)
	return &business
}

// TODO: Refactor to use the cache
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

func DeleteBusiness(userId, businessId uint) error {
	business, exists := businessCache.GetEntry(businessId)
	if !exists {
		return fmt.Errorf("Business with id: %d does not exist", businessId)
	}
	ei, err := api.CfxApi.DoRequest("DELETE", fmt.Sprintf("/business/%d", businessId), nil, nil)
	if err != nil {
		return err
	}
	if ei.Message != "" {
		return errors.New(ei.Message)
	}
	graylogger.Log("business:delete", fmt.Sprintf("%d removed a business", userId), "userId", userId, "business", business)
	return nil
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
	for i := range employees {
		employees[i].Role.Permissions = bitToPermList(employees[i].Role.PermMask)
	}
	return employees, err
}

func ChangeOwner(userId, businessId, newOwner uint) error {
	oldOwner := cfx_models.BusinessEmployee{}
	employee := cfx_models.BusinessEmployee{}
	err := db.CfxMariaDB.Client.First(&oldOwner, cfx_models.BusinessEmployee{BusinessId: businessId, IsOwner: true}).Error
	if err != nil {
		return err
	}
	err = db.CfxMariaDB.Client.First(&employee, cfx_models.BusinessEmployee{BusinessId: businessId, CitizenId: newOwner}).Error
	if err != nil {
		return err
	}
	if oldOwner.CitizenId == employee.CitizenId {
		return nil
	}
	oldOwner.IsOwner = false
	employee.IsOwner = true
	err = db.CfxMariaDB.Client.Save(&oldOwner).Error
	if err != nil {
		return err
	}
	err = db.CfxMariaDB.Client.Save(&employee).Error
	if err != nil {
		return err
	}
	graylogger.Log("business:updateOwner", fmt.Sprintf("%d has updated a business owner", userId), "businessId", businessId, "oldOwner", oldOwner.CitizenId, "newOwner", newOwner)

	cfxInput := struct {
		businessId uint
		newOwner   uint
	}{businessId: businessId, newOwner: newOwner}
	ai, err := api.CfxApi.DoRequest("POST", "/business/updateOwner", &cfxInput, nil)
	if err != nil {
		return err
	}
	if ai != nil && ai.Message != "" {
		return errors.New(ai.Message)
	}
	return nil
}
