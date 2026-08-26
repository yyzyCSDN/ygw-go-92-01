package model

type TrainState struct {
	ID      string
	Braking bool
	Power   float64
	Speed   float64
	Mass    float64
}

type BrakingReport struct {
	TrainID  string
	Power    float64
	Capacity float64
}

type Allocation struct {
	TrainID  string
	Capacity float64
	Power    float64
}

type DeviceSnapshot struct {
	ID       string
	Status   DeviceStatus
	State    AbsorberState
	Capacity float64
	Revision int
}

type SystemSnapshot struct {
	Revision int
	Devices  []DeviceSnapshot
}
