package realestate

import (
	"degrens/panel/internal/db"
	cfx_models "degrens/panel/internal/db/models/cfx"

	"github.com/aidenwallis/go-utils/utils"
)

func FetchLocationsForCitizenID(citizenID uint) ([]cfx_models.RealEstateLocation, error) {
	locationAccess := []cfx_models.RealEstateLocationAccess{}
	err := db.CfxMariaDB.Client.Where("cid = ?", citizenID).Find(&locationAccess).Error
	if err != nil {
		return nil, err
	}

	locations := []cfx_models.RealEstateLocation{}
	err = db.CfxMariaDB.Client.Preload("Access").Preload("Access.Character.Info").Where("id IN ?", utils.SliceMap(locationAccess, func(v cfx_models.RealEstateLocationAccess) uint {
		return v.LocationId
	})).Find(&locations).Error

	return locations, err
}
