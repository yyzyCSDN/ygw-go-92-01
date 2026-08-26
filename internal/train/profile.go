package train

import "regenbrake/internal/model"

type Profile struct {
	MaxPower   float64
	BaseMass   float64
	Efficiency float64
}

func (p Profile) BrakingPower(speed float64) float64 {
	if speed <= 0 {
		return 0
	}
	power := p.MaxPower * speed * p.Efficiency
	return power
}

func (p Profile) RequiredCapacity(power float64) float64 {
	if power <= 0 {
		return 0
	}
	return power * 0.8
}

func NewDefaultProfile() Profile {
	return Profile{MaxPower: 1000, BaseMass: 120, Efficiency: 0.9}
}

func (p Profile) Describe(state model.TrainState) map[string]float64 {
	power := p.BrakingPower(state.Speed)
	return map[string]float64{
		"power":    power,
		"capacity": p.RequiredCapacity(power),
	}
}
