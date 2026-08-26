package absorber

type EnergySnapshot struct {
	Absorbed   float64
	Discharged float64
	Cycles     int
}

type EnergyMeter struct {
	absorbed   float64
	discharged float64
	cycles     int
}

func NewEnergyMeter() *EnergyMeter {
	return &EnergyMeter{}
}

func (e *EnergyMeter) RecordAbsorb(energy float64) {
	e.absorbed += energy
	e.cycles++
}

func (e *EnergyMeter) RecordDischarge(energy float64) {
	e.discharged += energy
}

func (e *EnergyMeter) Absorbed() float64 {
	return e.absorbed
}

func (e *EnergyMeter) Discharged() float64 {
	return e.discharged
}

func (e *EnergyMeter) Cycles() int {
	return e.cycles
}

func (e *EnergyMeter) Efficiency() float64 {
	if e.absorbed <= 0 {
		return 0
	}
	return e.discharged / e.absorbed
}

func (e *EnergyMeter) Snapshot() EnergySnapshot {
	return EnergySnapshot{
		Absorbed:   e.absorbed,
		Discharged: e.discharged,
		Cycles:     e.cycles,
	}
}

func (e *EnergyMeter) Reset() {
	e.absorbed = 0
	e.discharged = 0
	e.cycles = 0
}
