package cfx_models

type VehicleState string

var (
	Parked    VehicleState = "parked"
	Out                    = "out"
	Impounded              = "impounded"
)

type PlayerVehicles struct {
	Vin       string       `json:"vin" gorm:"primaryKey"`
	CitizenId uint         `json:"citizenid" gorm:"column:cid"`
	Model     string       `json:"model"`
	Plate     string       `json:"plate"`
	FakePlate string       `json:"fakeplate"`
	State     VehicleState `json:"state"`
	GarageId  string       `json:"garageId"`
	Harness   int          `json:"harness"`
	Stance    string       `json:"stance"`
	Wax       int          `json:"wax"`
	Nos       int          `json:"nos"`
	Character Character    `gorm:"foreignKey:CitizenId"`
}
